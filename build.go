package main

import (
	"bufio"
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/cbroglie/mustache"
	"gopkg.in/yaml.v2"
)

// Template matches a base16 template file
type Template struct {
	Name      string
	Extension string
	Output    string
}

func readTemplates(path string) []*Template {
	osf, err := os.Open(path)
	if err != nil {
		panic(err)
	}
	defer osf.Close()
	reader := bufio.NewReader(osf)
	dec := yaml.NewDecoder(reader)
	configMap := make(map[string]*Template)
	err = dec.Decode(configMap)
	if err != nil {
		panic(err)
	}
	config := make([]*Template, 0, len(configMap))
	for k, v := range configMap {
		v.Name = k
		config = append(config, v)
	}
	return config
}

func (t *Template) Path() string {
	return filepath.Join("templates", t.Name+".mustache")
}

const schemesGit = "https://github.com/base16-project/base16-schemes.git"

func makeContext(scheme map[string]string) (map[string]string, error) {
	context := make(map[string]string)
	var err error
	getCheck := func(key string) string {
		value, ok := scheme[key]
		if !ok {
			err = fmt.Errorf("Missing %q in scheme", key)
		}
		return value
	}
	context["scheme-name"] = getCheck("scheme")
	context["scheme-author"] = getCheck("author")
	for i := 0; i < 16; i++ {
		context[fmt.Sprintf("base%02X-hex", i)] = strings.TrimLeft(getCheck(fmt.Sprintf("base%02X", i)), "#")
	}
	if err != nil {
		return nil, err
	}
	return context, nil
}

func yamlParseFile(filename string) (map[string]string, error) {
	f, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	defer f.Close()
	r := bufio.NewReader(f)
	dec := yaml.NewDecoder(r)
	m := make(map[string]string)
	err = dec.Decode(m)
	if err != nil {
		return nil, err
	}
	return m, nil
}

func renderTemplateToFile(parsedTemplate *mustache.Template, context map[string]string, filename string) error {
	str, err := parsedTemplate.Render(context)
	if err != nil {
		return err
	}
	err = ioutil.WriteFile(filename, []byte(str), 0644)
	return err
}

func downloadSchemes(dir string) {
	cmd := exec.Command("git", "clone", "--depth=1", schemesGit, dir)
	cmd.Dir = dir
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	err := cmd.Run()
	if err != nil {
		panic(fmt.Errorf("Error while cloning %s: %v", schemesGit, err))
	}
}

func main() {
	var err error

	templates := readTemplates("templates/config.yaml")
	parsedTemplates := make([]*mustache.Template, len(templates))
	for i, template := range templates {
		parsedTemplates[i], err = mustache.ParseFile(template.Path())
		if err != nil {
			panic(fmt.Errorf("Error while parsing %s: %v", template.Path(), err))
		}
	}

	var tempdir string

	if len(os.Args) == 1 {
		tempdir, err = ioutil.TempDir("", "base16.build.")
		if err != nil {
			panic(err)
		}
		fmt.Println(tempdir)
		downloadSchemes(tempdir)
	} else {
		tempdir = os.Args[1]
	}

	globPat := filepath.Join(tempdir, "*.yaml")
	sources, err := filepath.Glob(globPat)
	if err != nil {
		panic(fmt.Errorf("Error while globbing %s: %v", globPat, err))
	}
	for _, source := range sources {
		schemeMap, err := yamlParseFile(source)
		if err != nil {
			panic(fmt.Errorf("Error parsing yaml file %s: %v", source, err))
		}
		context, err := makeContext(schemeMap)
		if err != nil {
			panic(fmt.Errorf("Error making context for %s: %v", source, err))
		}
		for i, template := range templates {
			parsedTemplate := parsedTemplates[i]
			filename := filepath.Join(
				template.Output,
				fmt.Sprintf("base16-%s%s", strings.TrimSuffix(filepath.Base(source), ".yaml"), template.Extension))
			err = renderTemplateToFile(parsedTemplate, context, filename)
			if err != nil {
				panic(fmt.Errorf("Error writing %s", filename))
			}
		}
	}
	fmt.Printf("Generated %d schemes\n", len(sources))
}
