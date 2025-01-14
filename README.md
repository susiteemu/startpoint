# A CLI tool for managing and scripting http/restful requests: startpoint

| With Catppuccin Mocha (FTW!) |
|-------------------------|
| ![Dark mode FTW!](https://vhs.charm.sh/vhs-4a8Vw6refUQKYrVytUqSVi.gif)|

| Works with light mode too (Catppuccin Latte) |
|-------------------------|
| ![It can have a light mode too!](https://vhs.charm.sh/vhs-3mdLYEOUrLF39ADJ7WW9zZ.gif) |

<!-- toc -->

- [TL;DR Tell Me What and How](#tldr-tell-me-what-and-how)
- [Background and Motivation](#background-and-motivation)
- [My Plans for `startpoint`](#my-plans-for-startpoint)
- [Manual](#manual)
  * [Installation](#installation)
    + [On macOS](#on-macos)
  * [Commands](#commands)
  * [Requests TUI](#requests-tui)
    + [Features in *EDIT* mode](#features-in-edit-mode)
    + [Features in *SELECT* mode](#features-in-select-mode)
    + [Themes](#themes)
  * [Profiles TUI](#profiles-tui)
  * [Request Definitions](#request-definitions)
      - [A Note About Starlark and Runtime](#a-note-about-starlark-and-runtime)
    + [Different Requests](#different-requests)
      - [JSON](#json)
      - [Plain text](#plain-text)
      - [XML](#xml)
      - [Form data](#form-data)
      - [Multipart form data and Uploading files](#multipart-form-data-and-uploading-files)
      - [Downloading files](#downloading-files)
      - [URL/Query Parameters and Path Variables](#urlquery-parameters-and-path-variables)
    + [Chaining Requests](#chaining-requests)
    + [Templating Requests](#templating-requests)
  * [Profiles](#profiles)
  * [Importing](#importing)
  * [Configuration](#configuration)
  * [Examples](#examples)
- [TODO](#todo)

<!-- tocstop -->

## TL;DR Tell Me What and How

`startpoint` is a terminal-based application for managing and executing HTTP requests. It allows you to chain requests, use responses from previous requests, and automate complex workflows directly from the terminal.

To install `startpoint` see [Installation](#installation).

After installing you can run it with `startpoint`. This opens a TUI app to add/edit/manage requests. Add your own* and run. Done.

*)
You can test e.g. with this:

```yaml
url: https://httpbin.org/anything
method: POST
headers:
  Content-Type: 'application/json'
body: >
  {
    "id": 1,
    "name": "Jane"
  }
```

## Background and Motivation

When I got into the kind of web/API development where I had something else than SOAP to bite, that is RESTFUL APIs, I needed a proper tool to create, manage and run requests and view their responses. These APIs many times had some sort of authentication to them, which meant, excluding BASIC AUTH, some sort of request had to be made before the actual API call and then using the response from this authentication call with it. Other nice-to-have feature for me was profile/environment based configuration: to be able to run the same request but with different host etc without having the need to everytime rewrite the actual request itself or having many duplicates of it.

My first tool was Postman, which at that time was pretty neat: you had a nice user-interface, you could manage requests pretty easily and it had the concept of pre-request, which you could even define at the upper level that would apply to all of its sub-requests. Also it did have profiles. Unfortunately Postman became bloated with cloud/subscription stuff which are not what I seek.

Another nice find was Insomnia, which had a much simpler user-interface than what had become of bloated-Postman, it didn't flash all of the cloud/team collaboration stuff and subscription model to your face at the first meeting. I was pretty happy using it, even though it didn't have a direct support for pre-requests at the upper level. It did (and probably still does) support importing OpenAPI specificated requests so a bit of scripting and I could have what I needed. However, like Postman, also Insomnia fell into the cloud/subscription model.

Worth mentioning is Jetbrain's requests that come (at least) with IntelliJ IDEA. They are simple text files that contain definitions for a request. Easy to put into version control too. However I never got into them because a) I couldn't find a way to make pre-requests easily and b) they are run with IDEA which is a heavy-duty tool. Even though vast majority of my professional career has been spent coding Java, I wouldn't want to launch IDEA just to run http requests. Same reason I prefer separate database clients over what e.g. IntelliJ IDEA contains.

Also worth mentioning is Bruno. I heard about it some time ago, after I had already started implementing `startpoint`. It mentions its aim to revolutionize the status quo of current http client tools, it is offline-only, no cloud and its files are its own flavor of markup and easy to put into version control. Honestly, it sounds very good and is what I sought in Postman and Insomnia before. The only thing is...

Few years ago I got more and more interested in the terminal and running lightweight apps/tools/commands there instead of using apps that contain their own browsers with hundreds of megabytes (or even gigabytes) of RAM usage. I switched from (the excellent) Sublime Merge to `lazygit`, from Keepass XC to `gopass` and the latest from Visual Studio Code to `neovim`. Terminal, I have found, is a very powerful tool containing hundreds or thousands of heavy-duty commands to process any kind of data you need to. Having these other tools I use in my daily development work there keeps all nicely together in the same environment.

Then the idea got into my head. I had used and liked a lot `curl` and `httpie`, both tools for making http requests. The only caveat they had for me is that it would require some amount of scripting to achieve pre-requests and switchable profiles/environments. I started planning my own tool that supports of all my needs.

This is where `startpoint` comes in – a terminal tool that doesn't require a GUI, doesn't push cloud services, and is easy to integrate with version control.

## My Plans for `startpoint`

My plan and is to keep `startpoint` offline, cloud-free, subscription-free, ad-free, tracking-free and open-source. The definitions for requests, profiles and other metadata are kept as simple and as not-invented-here as possible. They can be put into the version control, zipped, Airdropped or even faxed if that is what you desire. I do not plan or promise to maintain it forever;  this is just my hobby project.

How do I see `startpoint`? In addition to what I mentioned above, I want it to be kind of fast, lightweight, pretty (I like my eyecandy), rather than constraining me to use its set of editors or other tools it would support me choosing my own, scriptable and extensible. I want to use standards and formats I and other people already know or if not, learning them should be beneficial outside of using this tool too.

As a more concrete list, at this point of time I have plans or have implemented:

- dotenv based profile support
- using yaml as request definitions
- using Starlark as scriptable request definitions
- using Lua as scriptable request definitions (coming later)
- support for using your own favorite editor in writing the request definitions
- TUI for managing both requests and profiles
- TUI for activating wanted profile and running requests
- printing requests pretty-formatted and colorized
- importing from OpenAPI specifications

## Manual

### Installation

There are several ways of installing `startpoint`.

#### On macOS

To install using Homebrew:

```bash
brew tap susiteemu/tap
brew install startpoint
```

To update using Homebrew:

```bash
brew update
brew upgrade startpoint
```

To install manually:

- Download tar.gz file from release
- Uncompress and either move to a location that is in your `$PATH` or run as a standalone from your desired location.

### Commands

There are few different commands.

```
❯ startpoint --help
Startpoint is a TUI app with which you can manage and run HTTP requests from your terminal. It offers a way for flexible chaining and scripting requests as well as defining them in a simple format.

Usage:
  startpoint [flags]
  startpoint [command]

Available Commands:
  help        Help about any command
  profiles    Start up a TUI application to manage profiles
  requests    Start up a TUI application to manage and run requests
  run         Run a http request from workspace

Flags:
      --config string      config file (default is a merge of $HOME/.startpoint.yaml and <workspace>/.startpoint.yaml)
      --help               Displays help
  -v, --version            version for startpoint
  -w, --workspace string   Workspace directory (default is current dir)
```

For each command there are some arguments/flags you can pass.

With `run` you can use flags/arguments to define which parts of the response to print.

```
❯ startpoint run --help
Run a http request from workspace

Usage:
  startpoint run [REQUEST NAME] [PROFILE NAME] [flags]

Flags:
      --no-body         Print no body
  -p, --plain           Print plain response without styling
      --print strings   Print WHAT
                        - 'h'   Print response headers
                        - 'b'   Print response body
                        - 't'   Print trace information

Global Flags:
      --config string      config file (default is a merge of $HOME/.startpoint.yaml and <workspace>/.startpoint.yaml)
      --help               Displays help
  -w, --workspace string   Workspace directory (default is current dir)
```

With `profiles` you can pass `workspace` and `config` file.

```
❯ startpoint profiles --help
Start up a TUI application to manage profiles

Usage:
  startpoint profiles [flags]

Global Flags:
      --config string      config file (default is a merge of $HOME/.startpoint.yaml and <workspace>/.startpoint.yaml)
      --help               Displays help
  -w, --workspace string   Workspace directory (default is current dir)
```

And the same with `run`.

```
❯ startpoint requests --help
Start up a TUI application to manage and run requests

Usage:
  startpoint requests [flags]

Global Flags:
      --config string      config file (default is a merge of $HOME/.startpoint.yaml and <workspace>/.startpoint.yaml)
      --help               Displays help
  -w, --workspace string   Workspace directory (default is current dir)
```

### Requests TUI

Requests TUI app has functionalities to add, edit, copy, remove, rename, preview and run requests and select active profile.

![Help](https://vhs.charm.sh/vhs-12avTKjVVhRTUyx8RIJaf3.gif)

TUI has two distinct *modes*: *SELECT* and *EDIT* modes.

- In *SELECT* mode you can change active profile, preview and run requests.
- In *EDIT* mode you can manage requests by adding, editing, copying, removing and renaming them. You can also preview them.

These are the keymappings for the *SELECT* mode:

```
 ↑/k      up             p preview             / filter
 ↓/j      down           r run
 →/l/pgdn next page      i edit mode
 ←/h/pgup prev page      a activate profile
 g/home   go to start
 G/end    go to end
```

And these for the *EDIT* mode:

```
 ↑/k      up             a   add          / filter
 ↓/j      down           e   edit
 →/l/pgdn next page      d   delete
 ←/h/pgup prev page      p   preview
 g/home   go to start    r   rename
 G/end    go to end      c   copy
                         esc view mode
```

You can quit the app with `q` or `ctrl+c`.

#### Features in *EDIT* mode

You can either *add* requests with `requests` TUI app or by creating `yaml` or `starlark` files directly with your favorite editor: it doesn't matter which way they are created. The app does not have any metadata mumbo-jumbo files to consider. At least in the beginning it is recommended to use the TUI app since it creates a template for you to use. To add a request with TUI app, press `a` when in *EDIT* mode. The app will ask if you want to add a `yaml` or `starlark` request. After this it will open your cup of editor (`$EDITOR`) and you can write the definition for it. After quitting the editor you can continue with the app.

![Adding a request](https://vhs.charm.sh/vhs-YBSuC3B9QKbyFVu5Q7Voa.gif)

*Editing* works the same: you can either do it directly by opening the request file with an editor or by opening the TUI app and starting the editing there.

*Deleting* means simply deleting the request file so it can also be done directly from file system or with the TUI app. The app does some checks before deleting (whether the request is defined as a previous request to other requests).

*Preview* opens the selected request to a syntax highlighted and scrollable view. Note that it shows the "raw" version of request and does not fill any template variables.

*Renaming* simply renames the request file so it can also be done through the file system. The app however has a ability to rename `prev_req` properties for other requests in case you are renaming a request that is used as a previous request to others.

*Copy* duplicates selected request, meaning you can use it as a base for a new request. Similar to many features above, you can do it either with the app or directly with your file system tools.

#### Features in *SELECT* mode

Requests can be *run* either with TUI app (pressing `r` when in *SELECT* mode) or directly with `run` command. At the moment there is no autocompletion when using `run` command so you would have to check the name of the request.

*Preview* is also available in this mode. It opens the selected request to a syntax highlighted and scrollable view. Note that it shows the "raw" version of request and does not fill any template variables.

With profile *activation* you tell the app to use variables from the profile and fill any possible template variables in the request. Read more about [Profiles](#profiles)

#### Themes

Currently there is no direct theming support, excluding syntax coloring which comes from [Chroma](https://github.com/alecthomas/chroma). However, you can assign colors to different configuration attributes in [Configuration](#configuration). As an example, [samples](samples) directory contains two different "themes": `.startpoint-catpuccing-latte.yaml` and `.startpoint-catppuccin-mocha.yaml`. The latter is also used as the default theme for the TUI apps.

### Profiles TUI

### Request Definitions

`startpoint` has different kinds of requests for varying needs: simple and more complex ones. Simple ones are defined with `yaml`, can be templated and run as a part of a request chain, but they can't use values from a previous request's response. Complex ones are scripted with `starlark`, can use values from both profile and previous request's response and use the bells and whistles of a programming language.

A request regardless of type has:

- *required* a name, which must be unique within the workspace. This is derived from the file name without extension.
- *required* an url to the resource
- *required* a request method: GET, POST, PUT, DELETE, PATCH, HEAD, OPTIONS
- *optional* the name of the previous request
- *optional* a list of headers
- *optional* a body
- *optional* a output path to response body (e.g. downloading a binary file)
- *optional* options defining things such as proxies, certificates e.g.

In addition there are some properties related to authentication which can be used as a shorthand.

A very simple example of request would be:

```yaml
# A request.yaml
url: https://httpbin.org/anything
method: GET
```

```python
# A request.star
url = "https://httpbin.org/anything"
method = "GET"
```

Both of these would perform a HTTP GET request url `https://httpbin.org/anything` and print the response.

To add some headers, you would do:

```yaml
# A request with headers.yaml
url: https://httpbin.org/anything
method: GET
headers:
  Accept: "application/json"
  X-Custom-Header: "Some custom value"
```

```python
# A request with headers.star
url = "https://httpbin.org/anything"
method = "GET"
headers = {
  "Accept": "application/json",
  "X-Custom-Header": "Some custom value"
}
```

And to add a body, you would do:

```yaml
# A request with body.yaml
url: https://httpbin.org/anything
method: GET
headers:
  Accept: "application/json"
  Content-Type: "application/json"
  X-Custom-Header: "Some custom value"
body: >
 {
   "id": 1,
   "name": "Jane"
 }
```

```python
# A request with body.star
url = "https://httpbin.org/anything"
method = "GET"
headers = {
  "Accept": "application/json",
  "Content-Type": "application/json",
  "X-Custom-Header": "Some custom value"
}
body = {
  "id": 1,
  "name": "Jane",
}
```

##### A Note About Starlark and Runtime

Since `starlark` is executable, values might be only resolved during runtime. This poses a challenge for cases when the app needs to know values for certain properties before running the script. These are:

- Requests TUI app displays request url and method in the listing. This is mainly to help you remember and know which request is which so it has only documentation purpose.
- When building request chain, previous request needs to be known before-hand.

For the first case, there are two possibilities:

- If your url and/or method is static (i.e. it does not change during runtime) you can just define it as is and the app will parse the value and show it.
- If your url and/or method changes during runtime, you can define a "multi-line comment block"* and add a "static version" there. Although what you see on the TUI would differ from what is the actual value when running the request, it could still be beneficial to see e.g. that your `method` is `GET or POST` instead of seeing `<blank>`.

*) An example of multi-line comment block

```python
"""
doc:url: http://localhost:8000/api/foo
doc:method: GET or POST
"""
```

For the second case, similar to the first one, you would also use the multi-line comment block:

```python
"""
prev_req: Some other request
"""
```

#### Different Requests

##### JSON

In `yaml` based requests you simply define the json body as a string:

```yaml
# ... other attributes...
body: '{"id": 1, "name": "Jane"}'

```

Or with "folded block scalar" (`>`):

```yaml
# ... other attributes...
body: >
  {
    "id": 1,
    "name": "Jane"
  }
```

In addition you (probably, depending on the endpoint you are using) need to include `Content-Type` header with the request:

```yaml
headers:
  Content-Type = "application/json"
```

In `starlark` based requests you can either define the json body as a string or as a dict. Which you choose depends on your needs: it is probably more convenient to use dictionaries when appending items dynamically.

As a string you would do:

```python
# ... other attributes...
body = '{"id": 1, "name": "Jane"}'
```

And as a dict you would do:

```python
# ... other attributes...
body = {
  "id": 1,
  "name": "Jane",
}
```

As with `yaml` requests, you probably want to add appropriate header:

```python
headers = {
  "Content-Type": "application/json"
}
```

##### Plain text

Plain text is plain and simple: pass it to all requests as a string. You also probably want to add appropriate headers.

With `yaml` based requests:

```yaml
headers:
  Content-Type: "text/plain"
body: "Plain text body"
```

And with `starlark` based requests:

```python
headers = {
  "Content-Type": "text/plain"
}
body = "Plain text body"
```

##### XML

XML formatted body is passed simply as a string. You probably also need to include correct header.

With `yaml` based requests you could do:

```yaml
headers:
  Content-Type: "application/xml"
body: >
  <root>
    <id>1</id>
    <name>Jane</name>
  </root>
```

And with `starlark` based requests you could do:

```python
headers = {
  "Content-Type": "application/xml"
}
body = """
  <root>
    <id>1</id>
    <name>Jane</name>
  </root>
"""
```

##### Form data

You can send form data by adding either `application/x-www-form-urlencoded` or `multipart/form-data` as `Content-Type` header (see more about [Multipart form data](#multipart-form-data-and-uploading-files)) and defining the body as a map/dict.

With `yaml` based requests:

```yaml
method: POST
headers:
  Content-Type: 'application/x-www-form-urlencoded'
body:
  field1: val1
  field2: val2
  field3: val3
```

With `starlark` based requests:

```python
method = "POST"
headers = {
  "Content-Type": "application/x-www-form-urlencoded"
}
body = {
  "field1": "val1",
  "field2": "val2",
  "field3": "val3",
}
```

##### Multipart form data and Uploading files

Sending multipart form data and uploading files is possible by adding a) appropriate header and b) defining the body as a map/dict.

If you want to upload a file, the map/dict entry should begin with `@` followed by the path to the file.

With `yaml` based requests you can send multipart form data like this:

```yaml
method: POST
headers:
  Content-Type: 'multipart/form-data'
body:
  title: 'Image title'
  file: '@resources/Image.png'
```

And with `starlark` based requests like this:

```python
method = "POST"
headers = {
  "Content-Type": "multipart/form-data"
}
body = {
  "title": "Image title",
  "file": "@resources/Image.png",
}
```

##### Downloading files

When you want to download the response instead of printing it, which would be sensible especially when response is a binary file, you define `output` property and point it to a file you want the response be saved to.

With `yaml` based requests you would do:

```yaml
url: "http://localhost:8000/download"
output: "/path/to/some/file.png"
```

With `starlark` based requests you would do:

```python
url = "http://localhost:8000/download"
output = "/path/to/some/file.png"
```

##### URL/Query Parameters and Path Variables

Currently there is no special property for passing query parameters and path variables. Instead you would simply add them to the url property. Check [templating](#templating-requests) of how to pass values from environment. With `starlark` based requests you can also assign values dynamically.

Query parameters with `yaml` based requests:

```yaml
url: "http://localhost:8000?arg1=val1&arg2=val2"
```

Query parameters with `starlark` based requests:

```python
url = "http://localhost:8000?arg1=val1&arg2=val2"
```

Path variables with `yaml` based requests:

```yaml
url: "http://localhost:8000/some-path-var/123"
```

Path variables with `starlark` based requests:

```python
url = "http://localhost:8000/some-path-var/123"
```

#### Chaining Requests

At times it is useful to run a request before another, e.g. when using a API that has a authentication scheme requiring to pass a token. Each request, regardless of being "simple" or "complex" has a property `prev_req` that can be used to point to a another request. When used with "simple" (`yaml` based) requests you can't use values from the previous response but you can nevertheless chain them if need be. The real benefit comes when using "complex" (`starlark` based) requests: you can take values from previous response's headers and body, build logic upon them and pass them to the current request.

With `yaml` based requests you can define previous request like so:

```yaml
# ... other attributes...
prev_req: Some other request
```

With `starlark` based requests you can do it like so:

```python
"""
prev_req: Some other request
"""

```

An example illustrates how to authenticate to oauth2 endpoint.

This is the `User details.star` request. It wants to perform a `GET` request to `/auth/oauth2/users/me` endpoint that returns data about the user. The endpoint is protected with oauth2 which basically means you have to pass a header `Authorization` with the value of `Bearer + <access token>` in order to authenticate. This request has a previous request defined `Token` from which it gets `prevResponse` dictionary/map. Using this map it accesses the previous response's body and from the body `access_token` attribute.

```python
"""
prev_req: Token
"""
url = "http://localhost:8000/auth/oauth2/users/me/"
method = "GET"
auth = "Bearer " + prevResponse["body"]["access_token"]
headers = { "Authorization": auth }
```

The `Token.yaml` request is as follows. It performs a `POST` request to `/auth/oauth2/token` with form data holding the user credentials. Note that form fields depend on which authentication flow is used. Note also, that it is not advised to put sensitive values such as passwords directly to requests but use the [templating](#templating-yaml-requests) mechanism.

```yaml
url: http://localhost:8000/auth/oauth2/token
method: POST
headers:
  Content-Type: 'application/x-www-form-urlencoded'
  Accept: 'application/x-www-form-urlencoded, application/json'
body:
  username: 'johndoe'
  password: 'secret'
  grant_type: 'password'
```

#### Templating Requests

It is possible and often useful to template request values: this way you can use the same request definition in different profiles/environments.

To template a value, just use `{value_name}` syntax. The `value_name` refers to variable defined in the profile file. You can template:

- the url of the request
- header names and/or header values
- parts of body

An example: you want to perform a `GET` request to an endpoint `/foo`. You have multiple environments you want to ultimately test/try. You define the request as follows:

```yaml
# Using templates.yaml
url: {domain}/foo
method: GET
```

```python
# Using templates.star
url = "{domain}/foo"
method = "GET"
```

Your profiles files could then look like this:

```bash
# .env
domain=http://localhost:8000
```

```bash
# .env.test
domain=https://yourtestdomain.com
```

```bash
# .env.prod
domain=https://yourdomain.com
```

Now, when you run your request in the `default` (`.env`) profile, your url would be `http://localhost:8000/foo`. In `test` it would be `https://yourtestdmain.com/foo` and in `prod` `https://yourdomain.com/foo`.

### Profiles

Profiles are a way to run requests with different groups of variables. You can e.g. have one profile for your local environment holding request urls such as `http://localhost:8080` etc and one for your prod environment having its own urls. When you define profiles and variables inside them, you can use these variables in requests allowing you to avoid hard-coding values and reusing same request definitions on different situations and needs.

Profiles are based on dotenv, meaning you can define a "default" profile by creating a file called `.env` and filling it with variables like:

```bash
domain=http://localhost:8080
```

To create different profiles, you would add more files: `.env.test`, `.env.prod` etc.

Profiles are merged and the priority goes (from lowest to highest):

- `.env`
- `.env.local`
- `.env.<selected profile>`
- `.env.<selected profile>.local`

Profiles with suffix `.local` are meant to hold sensitive values such as passwords. Whereas you can put other files to version control, it is recommended that you keep `.local` files out of it.

### Importing

You can import requests and profiles (a workspace) from OpenAPI specifications. Currently only version 3 is supported. To import, you can use the `import` command:

```
❯ startpoint import --help
Import workspace from OpenAPI Spec v3

Usage:
  startpoint import [flags]

Flags:
  -p, --path string   OpenAPI Spec v3 location (filepath or url)

Global Flags:
      --config string      config file (default is a merge of $HOME/.startpoint.yaml and <workspace>/.startpoint.yaml)
      --help               Displays help
  -w, --workspace string   Workspace directory (default is current dir)
```

### Configuration

Not to confuse with profiles, which are used to define variables to requests, the configuration is metadata used to change the look and behaviour of the application.

The configuration allows you to customize the look of the ui, set some properties for the HTTP client (e.g. proxy settings), change logging levels, define whether a response parts will print etc.

Configuration can come from multiple sources:

- from a config file in user home directory
- from a config file in the used workspace
- from an explicitly defined config file
- from environment variables
- from the request that is being run

Whenever there are more than one configuration source, the final effective configuration will be a merge of all of them, the base being the one with broadest scope.

If you do not explicitly define a config file, the lookup and merge order will be:

- take base from the config file in user home directory (if exists)
- merge with configuration coming from the config file in workspace (if exists)
- merge with configuration coming from environment variables
- when running a request, merge with configuration coming from this request

If you explicitly define a config file, the corresponding order will be:

- take base from defined config file
- merge with configuration coming from environment variables
- when running a request, merge with configuration coming this request

The configuration is given in `yaml` format (excluding those coming from environment). An example of configuration file is [here](./samples/.startpoint-all-configurations.yaml).

All configuration values are:

| Key | Default Value | Description | Scope |
| ------------- | -------------- | -------------- | -------------- |
| theme.syntax | catppuccin-mocha | Sets syntax coloring for [Chroma](https://github.com/alecthomas/chroma). See [all available styles](https://github.com/alecthomas/chroma/tree/master/styles). | Global |
| theme.bgColor | ![#1e1e2e](https://via.placeholder.com/15/1e1e2e/000000.png?text=+) `#1e1e2e` | Background color for the TUI apps | Global |
| theme.primaryTextFgColor | ![#cdd6f4](https://via.placeholder.com/15/cdd6f4/000000.png?text=+) `#cdd6f4` | Primary text foreground color | Global |
| theme.secondaryTextFgColor | ![#bac2de](https://via.placeholder.com/15/bac2de/000000.png?text=+) `#bac2de`  | Secondary text foreground color | Global |
| theme.titleFgColor |  ![#1e1e2e](https://via.placeholder.com/15/1e1e2e/000000.png?text=+) `#1e1e2e`  | App title foreground color | Global |
| theme.titleBgColor | ![#a6e3a1](https://via.placeholder.com/15/a6e3a1/000000.png?text=+) `#a6e3a1`  | App title background color | Global |
| theme.borderFgColor |  ![#cdd6f4](https://via.placeholder.com/15/cdd6f4/000000.png?text=+) `#cdd6f4` | Border foreground color | Global |
| theme.whitespaceFgColor | ![#313244](https://via.placeholder.com/15/313244/000000.png?text=+) `#313244`  | Foreground color for the whitespace (shown as a background for dialogs/prompts) | Global |
| theme.errorFgColor | ![#f38ba8](https://via.placeholder.com/15/f38ba8/000000.png?text=+) `#f38ba8`  | Foreground color for errors | Global |
| theme.statusbar.primaryBgColor | ![#11111b](https://via.placeholder.com/15/11111b/000000.png?text=+) `#11111b`   | Background color for the primary section of the statusbar (e.g. that displays messages) | Global |
| theme.statusbar.primaryFgColor |  ![#cdd6f4](https://via.placeholder.com/15/cdd6f4/000000.png?text=+) `#cdd6f4`  | Foreground color for the primary section of the statusbar (e.g. that displays messages) | Global |
| theme.statusbar.secondaryFgColor |  ![#1e1e2e](https://via.placeholder.com/15/1e1e2e/000000.png?text=+) `#1e1e2e`  | Foreground color for "secondary" content such as other, colored, sections | Global |
| theme.statusbar.modePrimaryBgColor | ![#f9e2af](https://via.placeholder.com/15/f9e2af/000000.png?text=+) `#f9e2af`   | Background color for the "primary" (e.g. in requests TUI the *SELECT*) mode section of the statusbar| Global |
| theme.statusbar.modeSecondaryBgColor |  ![#a6e3a1](https://via.placeholder.com/15/a6e3a1/000000.png?text=+) `#a6e3a1`   | Background color for the "secondary" (e.g. in requests TUI the *EDIT*) mode section of the statusbar| Global |
| theme.statusbar.thirdColBgColor | ![#94e2d5](https://via.placeholder.com/15/94e2d5/000000.png?text=+) `#94e2d5` | Background color for the third section of the statusbar | Global |
| theme.statusbar.fourthColBgColor | ![#89b4fa](https://via.placeholder.com/15/89b4fa/000000.png?text=+) `#89b4fa`  | Background color for the fourth section of the statusbar | Global |
| theme.httpMethods.textFgColor |  ![#1e1e2e](https://via.placeholder.com/15/1e1e2e/000000.png?text=+) `#1e1e2e`  | In list items, the foreground color for the label showing request's http method | Global |
| theme.httpMethods.defaultBgColor |  ![#cdd6f4](https://via.placeholder.com/15/cdd6f4/000000.png?text=+) `#cdd6f4`  | In list items, the default background color for the label showing request's http method | Global |
| theme.httpMethods.getBgColor |  ![#89b4fa](https://via.placeholder.com/15/89b4fa/000000.png?text=+) `#89b4fa`   | In list items, the background color for the label showing request's http method when the method is `GET` | Global |
| theme.httpMethods.postBgColor |  ![#a6e3a1](https://via.placeholder.com/15/a6e3a1/000000.png?text=+) `#a6e3a1`   | In list items, the background color for the label showing request's http method when the method is `POST` | Global |
| theme.httpMethods.putBgColor |  ![#f9e2af](https://via.placeholder.com/15/f9e2af/000000.png?text=+) `#f9e2af`    | In list items, the background color for the label showing request's http method when the method is `PUT` | Global |
| theme.httpMethods.deleteBgColor | ![#f38ba8](https://via.placeholder.com/15/f38ba8/000000.png?text=+) `#f38ba8`   | In list items, the background color for the label showing request's http method when the method is `DELETE` | Global |
| theme.httpMethods.patchBgColor | ![#94e2d5](https://via.placeholder.com/15/94e2d5/000000.png?text=+) `#94e2d5`  | In list items, the background color for the label showing request's http method when the method is `PATCH` | Global |
| theme.httpMethods.optionsBgColor |  | In list items, the background color for the label showing request's http method when the method is `OPTIONS` | Global |
| theme.urlFgColor |  ![#89b4fa](https://via.placeholder.com/15/89b4fa/000000.png?text=+) `#89b4fa`   | In list items, the foreground color for request's URL | Global |
| theme.urlBgColor |  |  In list items, the background color for request's URL | Global |
| theme.urlTemplatedSectionFgColor |  ![#f9e2af](https://via.placeholder.com/15/f9e2af/000000.png?text=+) `#f9e2af`    | In list items, the foreground color for the templated section of the request's URL | Global |
| theme.urlTemplatedSectionBgColor |  | In list items, the background color for the templated section of the request's URL | Global |
| theme.urlUnfilledTemplatedSectionFgColor | ![#f38ba8](https://via.placeholder.com/15/f38ba8/000000.png?text=+) `#f38ba8`   |In list items, the foreground color for the templated section of the request's URL when there is no environment/profile value to match the templated variable | Global |
| theme.urlUnfilledTemplatedSectionBgColor |  |  In list items, the foreground color for the templated section of the request's URL when there is no environment/profile value to match the templated variable | Global |
| theme.response.status200FgColor |  ![#a6e3a1](https://via.placeholder.com/15/a6e3a1/000000.png?text=+) `#a6e3a1`   | Foreground color for the response's 2xx status | Global |
| theme.response.status300FgColor |  ![#f9e2af](https://via.placeholder.com/15/f9e2af/000000.png?text=+) `#f9e2af`    | Foreground color for the response's 3xx status | Global |
| theme.response.status400FgColor |  ![#f38ba8](https://via.placeholder.com/15/f38ba8/000000.png?text=+) `#f38ba8`   | Foreground color for the response's 4xx status | Global |
| theme.response.status500FgColor |  ![#f38ba8](https://via.placeholder.com/15/f38ba8/000000.png?text=+) `#f38ba8`   | Foreground color for the response's 5xx status | Global |
| theme.response.protoFgColor |  ![#89b4fa](https://via.placeholder.com/15/89b4fa/000000.png?text=+) `#89b4fa`   | Foreground color for the response's proto part | Global |
| theme.response.headerFgColor |  ![#89b4fa](https://via.placeholder.com/15/89b4fa/000000.png?text=+) `#89b4fa`   | Foreground color for the response's header names | Global |
| printer.pretty | `true`| Pretty print responses | Global, request |
| printer.formatter | terminal16m | Which formatter to use with Chroma | Global |
| editor | `$EDITOR` | Which editor to use for creating/editing requests and profiles | Global |
| debug | `false` | Enable debug logging | Global |
| httpClient.debug | `false` | Enable debug logging for the http client | Global, request |
| httpClient.enableTraceInfo | `false` | Include and print traceinfo with the response | Global, request |
| httpClient.insecure | `false` | Disable security check for https | Global, request |
| httpClient.proxyUrl | | Set proxy | Global, request |
| httpClient.timeoutSeconds | | Set timeout in seconds | Global, request |
| httpClient.clientCertificates[].certFile | | Array of certFile and keyFile pairs; certFile contains path to the public key file | Global, request |
| httpClient.clientCertificates[].keyFile | | Array of certFile and keyFile pairs; keyFile contains path to the private key file | Global, request |
| httpClient.rootCertificates[] | | Array of paths to custom root certificates | Global, request |

### Examples

Request without body or headers:

```yaml
# Yaml
url: https://httpbin.org/anything
method: GET
```

Request with body and headers:

```yaml
# Yaml
url: https://httpbin.org/anything
method: POST
headers:
  Content-Type: "application/json"
body: >
  {
    "id": 1,
    "name": "Jane"
  }
```

Request with formdata:

```yaml
# Yaml
url: https://httpbin.org/anything
method: POST
headers:
  Content-Type: 'application/x-www-form-urlencoded'
body:
  field1: val1
  field2: val2
  field3: val3
```

Request with basic auth:

```yaml
# Yaml
url: https://httpbin.org/basic-auth/someuser/somepassword
method: GET
auth:
  basic_auth:
    user: someuser
    password: somepassword
```

Request with bearer token:

```yaml
# Yaml
url: https://example.com/auth-with-bearer-token
method: GET
auth:
  bearer_token: some-token
```

Request with file output:

```yaml
# Yaml
url: https://example.com/somefile.pdf
method: GET
output: ./somefile.pdf
```

Request with options:

```yaml
# Yaml
url: https://example.com/somefile.pdf
method: GET
options:
  print: false
  debug: true
  enableTrace: true
  printRequest: false
```

## TODO

There are things still in progress and planned for some later date.

- [ ] Add Lua based requests v.1.2
- [ ] Preview, when a profile is selected, could auto-fill variables (but also show there's a variable; nvim "virtualtext" like?)
