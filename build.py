#!/usr/bin/env python3

import io
import os
import subprocess
import tempfile
import glob
import urllib.request

import pystache
import yaml


def make_context(scheme):
    context = {}
    context['scheme-name'] = scheme['scheme']
    context['scheme-author'] = scheme['author']
    for base in range(16):
        context[f'base{base:02X}-hex'] = scheme[f'base{base:02X}']
    return context

with open('templates/config.yaml') as file:
    config = yaml.safe_load(file)

templates = []
for path, details in config.items():
    with open(os.path.join('templates', path + '.mustache')) as file:
        template_source = file.read()
    parsed = pystache.parse(template_source)
    templates.append((parsed, details['extension'], details['output']))

with urllib.request.urlopen('https://github.com/chriskempson/base16-schemes-source/raw/master/list.yaml') as fileb:
    file = io.TextIOWrapper(fileb)
    repos = yaml.safe_load(file)

with tempfile.TemporaryDirectory() as tempdir:
    for repo, url in repos.items():
        subprocess.run(
            ['git', 'clone', '--depth=1', url, repo],
            cwd=tempdir,
            check=True)
    for schemepath in glob.iglob(os.path.join(tempdir, '**', '*.yaml')):
        with open(schemepath) as schemefile:
            scheme = yaml.safe_load(schemefile)
        context = make_context(scheme)
        for parsed, extension, outdir in templates:
            filename = os.path.join(
                outdir,
                'base16-' + os.path.basename(schemepath[:-5]) + extension
            )
            print('render', filename)
            with open(filename, 'w') as file:
                file.write(
                    pystache.render(parsed, context)
                )
