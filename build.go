package main

import (
	"bufio"
	"fmt"
	"io/ioutil"
	"net/http"
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

const schemesURL = "https://github.com/chriskempson/base16-schemes-source/raw/master/list.yaml"

// Scheme matches a base16 scheme source
type Scheme struct {
	Name string
	Repo string
}

func readSchemes(url string) []*Scheme {
	resp, err := http.Get(url)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		panic(fmt.Errorf("unexpected status code %d from %s", resp.StatusCode, url))
	}
	dec := yaml.NewDecoder(resp.Body)
	schemeMap := make(map[string]string)
	err = dec.Decode(schemeMap)
	if err != nil {
		panic(err)
	}
	schemes := make([]*Scheme, 0, len(schemeMap))
	for k, v := range schemeMap {
		schemes = append(schemes, &Scheme{k, v})
	}
	return schemes
}

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
		context[fmt.Sprintf("base%02X-hex", i)] = getCheck(fmt.Sprintf("base%02X", i))
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

func downloadSchemes(dir string, schemes []*Scheme) {
	doScheme := make(chan *Scheme, len(schemes))
	schemeDone := make(chan *Scheme)
	for i := 0; i < 4; i++ {
		go func() {
			for scheme := range doScheme {
				cmd := exec.Command("git", "clone", "--quiet", "--depth=1", scheme.Repo, scheme.Name)
				cmd.Dir = dir
				cmd.Stdout = os.Stdout
				cmd.Stderr = os.Stderr
				err := cmd.Run()
				if err != nil {
					panic(fmt.Errorf("Error while cloning %s: %v", scheme.Repo, err))
				}
				schemeDone <- scheme
			}
		}()
	}

	for _, scheme := range schemes {
		doScheme <- scheme
	}
	close(doScheme)

	nDone := 0
	for nDone < len(schemes) {
		nDone++
		fmt.Printf("(%d/%d) %s\n", nDone, len(schemes), (<-schemeDone).Name)
	}
}

func main() {
	var err error

	templates := readTemplates("templates/config.yaml")
	fmt.Println(templates)
	parsedTemplates := make([]*mustache.Template, len(templates))
	for i, template := range templates {
		parsedTemplates[i], err = mustache.ParseFile(template.Path())
		if err != nil {
			panic(fmt.Errorf("Error while parsing %s: %v", template.Path(), err))
		}
	}

	var tempdir string

	if len(os.Args) == 1 {
		schemes := readSchemes(schemesURL)
		tempdir, err = ioutil.TempDir("", "base16.build.")
		if err != nil {
			panic(err)
		}
		fmt.Println(tempdir)
		downloadSchemes(tempdir, schemes)
	} else {
		tempdir = os.Args[1]
	}

	globPat := filepath.Join(tempdir, "**", "*.yaml")
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
}
