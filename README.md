# A cli tool for managing and scripting http/restful requests: startpoint

```
###############################
#                             #
#       WORK IN PROGRESS      #
#                             #
###############################
```

<!-- toc -->

- [TL;DR Tell Me What and How](#tldr-tell-me-what-and-how)
- [Background and Motivation](#background-and-motivation)
- [My Plans for `startpoint`](#my-plans-for-startpoint)
- [Manual](#manual)
  * [Quickstart](#quickstart)
  * [Installing](#installing)
  * [Adding and Running Requests](#adding-and-running-requests)
  * [Commands](#commands)
  * [Request Composition](#request-composition)
    + [Different Content Types](#different-content-types)
    + [Chaining Requests](#chaining-requests)
    + [Templating `yaml` Requests](#templating-yaml-requests)
    + [Advanced `starlark` Requests](#advanced-starlark-requests)
  * [Profiles](#profiles)
  * [Configuration](#configuration)
  * [Examples](#examples)

<!-- tocstop -->

## TL;DR Tell Me What and How

`startpoint` is a terminal ui / terminal app for managing and running HTTP requests. You can chain requests, use values from previous request's response and do lots of other kind script magic.

To install `startpoint` do XXX

After installing you can run it with `startpoint requests`. This opens a ui to add/edit/manage requests. Add your own* and run. Done.

*)
You can test e.g. with this:

```
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

## Manual

### Quickstart



### Installing


### Adding and Running Requests

You can either add requests with `requests` tui app or by creating `yaml` or `starlark` files directly with your favorite editor: it doesn't matter which way they are created. The app does not have any metadata mumbo-jumbo files to consider. At least in the beginning it is recommended to use the tui app since it creates a template for you to use.

Requests can be run either with tui app or directly with `run` command. At the moment there is no autocompletion when using `run` command so you would have to check the name of the request.


### Commands

There are few different commands.

```
❯ startpoint --help
Startpoint is a tui app with which you can manage and run HTTP requests from your terminal. It offers a way for flexible chaining and scripting requests as well as defining them in a simple format.

Usage:
  startpoint [flags]
  startpoint [command]

Available Commands:
  help        Help about any command
  profiles    Start up a tui application to manage profiles
  requests    Start up a tui application to manage and run requests
  run         Run a http request from workspace
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
Start up a tui application to manage profiles

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
Start up a tui application to manage and run requests

Usage:
  startpoint requests [flags]

Global Flags:
      --config string      config file (default is a merge of $HOME/.startpoint.yaml and <workspace>/.startpoint.yaml)
      --help               Displays help
  -w, --workspace string   Workspace directory (default is current dir)
```


### Request Composition

`startpoint` has different kinds of requests for varying needs: simple and more complex ones. Simple ones are defined with `yaml`, can be templated and run as a part of a request chain, but they can't use values from a previous request's response. Complex ones are scripted with `starlark`, can use values from both profile and previous request's response and use the bells and whistles of a programming language.

A request regardless of type has:
* *required* a name, which must be unique within the workspace. This is derived from the file name without extension.
* *required* an url to the resource
* *required* a request method: GET, POST, PUT, DELETE, PATCH, HEAD, OPTIONS
* *optional* the name of the previous request
* *optional* a list of headers
* *optional* a body
* *optional* a output path to response body (e.g. downloading a binary file)
* *optional* options defining things such as proxies, certificates e.g.

In addition there are some properties related to authentication which can be used as a shorthand.

A very simple example of request would be:

```
# A request.yaml
url: https://httpbin.org/anything
method: GET

### Or with starlark

# A request.star
url = "https://httpbin.org/anything"
method = "GET"
```
Both of these would perform a HTTP GET request url `https://httpbin.org/anything` and print the response.

To add some headers, you would do:

```
# A request with headers.yaml
url: https://httpbin.org/anything
method: GET
headers:
  Accept: "application/json"
  X-Custom-Header: "Some custom value"

### Or with starlark

# A request with headers.star
url = "https://httpbin.org/anything"
method = "GET"
headers: {
  "Accept": "application/json",
  "X-Custom-Header": "Some custom value"
}
```
And to add a body, you would do:

```
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

### Or with starlark

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

#### Different Content Types

#### Chaining Requests

#### Templating `yaml` Requests

#### Advanced `starlark` Requests

### Profiles

Profiles are a way to run requests with different groups of variables. You can e.g. have one profile for your local environment holding request urls such as `http://localhost:8080` etc and one for your prod environment having its own urls. When you define profiles and variables inside them, you can use these variables in requests allowing you to avoid hard-coding values and reusing same request definitions on different situations and needs.

Profiles are based on dotenv, meaning you can define a "default" profile by creating a file called `.env` and filling it with variables like:
```
DOMAIN=http://localhost:8080
```
To create different profiles, you would add more files: `.env.test`, `.env.prod` etc.

Profiles are merged and the priority goes (from lowest to highest):
* `.env`
* `.env.local`
* `.env.<selected profile>`
* `.env.<selected profile>.local`

Profiles with suffix `.local` are meant to hold sensitive values such as passwords. Whereas you can put other files to version control, it is recommended that you keep `.local` files out of it.

### Configuration

Not to confuse with profiles, which are used to define variables to requests, the configuration is metadata used to change the look and behaviour of the application.

The configuration allows you to customize the look of the ui, set some properties for the HTTP client (e.g. proxy settings), change logging levels, define whether a response will print etc.

Configuration can come from multiple sources:
* from a config file in user home directory
* from a config file in the used workspace
* from a explicitly defined config file
* from environment variables
* from the request that is being run

Whenever there are more than one configuration source, the final effective configuration will be a merge of all of them, the base being the one with broadest scope.

If you do not explicitly define a config file, the lookup and merge order will be:
* take base from the config file in user home directory (if exists)
* merge with configuration coming from the config file in workspace (if exists)
* merge with configuration coming from environment variables
* when running a request, merge with configuration coming from this request

If you explicitly define a config file, the corresponding order will be:
* take base from defined config file
* merge with configuration coming from environment variables
* when running a request, merge with configuration coming the this request

All configuration values are:
```
theme:
  syntax: catppuccin-mocha
  bgColor: '#1e1e2e'
  primaryTextFgColor: '#cdd6f4'
  secondaryTextFgColor: '#bac2de'
  titleFgColor: '#1e1e2e'
  titleBgColor: '#a6e3a1'
  borderFgColor: '#cdd6f4'
  whitespaceFgColor: '#313244'
  statusbar:
    primaryBgColor: '#11111b'
    primaryFgColor: '#cdd6f4'
    secondaryFgColor: '#1e1e2e'
    modePrimaryBgColor: '#f9e2af'
    modeSecondaryBgColor: '#a6e3a1'
    thirdColBgColor: '#94e2d5'
    fourthColBgColor: '#89b4fa'
  httpMethods:
    textFgColor: '#1e1e2e'
    defaultBgColor: '#cdd6f4'
    getBgColor: '#89b4fa'
    postBgColor: '#a6e3a1'
    putBgColor: '#f9e2af'
    deleteBgColor: '#f38ba8'
    patchBgColor: '#94e2d5'
    optionsBgColor: ''
  urlFgColor: '#89b4fa'
  urlBgColor: ''
  urlTemplatedSectionFgColor: '#f9e2af'
  urlTemplatedSectionBgColor: ''
  urlUnfilledTemplatedSectionFgColor: '#f38ba8'
  urlUnfilledTemplatedSectionBgColor: ''
printer:
  response:
    formatter: terminal16m
editor: /opt/homebrew/bin/nvim
httpClient:
  debug: true
  enableTraceInfo: true
  insecure: false
  proxyUrl: https://some-proxy.com
  timeoutSeconds: 60
  clientCertificates:
    - certFile: /path/to/client.pem
      keyFile: /path/to/client.key
  rootCertificates:
    - /path/to/rootcert.pem
    - /path/to/another/rootcert.pem
```



### Examples

Request without body or headers:
```
# Yaml
url: https://httpbin.org/anything
method: GET
```

Request with body and headers:
```
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
```
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

Request with basic auth (NOTE: not implemented yet):
```
# Yaml
url: https://httpbin.org/basic-auth/someuser/somepassword
method: GET
auth:
  basic:
    user: someuser
    password: somepassword
```

Request with bearer token (NOTE: not implemented yet):
```
# Yaml
url: https://example.com/auth-with-bearer-token
method: GET
auth:
  bearer_token: some-token
```

Request with file output:
```
# Yaml
url: https://example.com/somefile.pdf
method: GET
output: ./somefile.pdf
```

Request with options:
```
# Yaml
url: https://example.com/somefile.pdf
method: GET
options:
  print: false
  debug: true
  enableTrace: true
```

## TODO

There are things still in progress and planned for some later date.

- [x] Tui for profiles v.1.0
- [w] Add logging v.1.0
- [x] Make configurable things configurable v.1.0
- [w] Add README.md v.1.0
- [ ] Setup ci/cd v.1.0
- [w] Add support/test different payloads
- [w] Create server to test requests/check if httpbin can be used effectively
- [ ] Add check for deleting request whether is breaks other requests v.1.0
- [ ] Fix bug in copy/edit/rename during filtering: before setting item, cancel filtering, maybe it fixes it? And could be logical: why continue filtering after that v.1.0
- [w] Http client settings: proxies, timeouts, trace logging, ... v1.0 (partly at least)
- [ ] Renaming request: also rename prev_req for other requests v.1.0
- [ ] Import from openspec v.1.1
- [ ] Shell completions: https://github.com/spf13/cobra/blob/main/site/content/completions/_index.md v.1.1
- [ ] Env "local"
- [ ] Add basicAuth and authToken so no need to include them into headers? v.1.1
- [ ] Add print request? v.1.1
- [ ] Add continueOnPrevRequestStatus etc v.1.1
- [w] Make failures (running request fails) more pretty/informative
