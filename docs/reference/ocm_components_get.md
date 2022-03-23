## ocm components get

get component version

### Synopsis


Get lists all component versions specified, if only a component is specified
all versions are listed.

If the repository option is specified, the given names are interpreted
relative to the specified repository using the syntax

<center><code>&lt;component>[:&lt;version>]</code></center>

If no <code>repo</code> option is specified the given names are interpreted 
as located OCM component version references:

<center><code>[&lt;repo type>::]&lt;host>[:&lt;port>][/&lt;base path>]//&lt;component>[:&lt;version>]</code></center>

Additionally there is a variant to denote common transport archives
and general repository specifications

<center><code>[&lt;repo type>::]&lt;filepath>|&lt;spec json>[//&lt;component>[:&lt;version>]]</code></center>

The <code>--repo</code> option takes an OCM repository specification:

<center><code>[&lt;repo type>::]&lt;configured name>|&lt;file path>|&lt;spec json></code></center>

For the *Common Transport Format* the types <code>directory</code>,
<code>tar</code> or <code>tgz</code> is possible.

Using the JSON variant any repository type supported by the 
linked library can be used:

Dedicated OCM repository types:
- `ComponentArchive`

OCI Repository types (using standard component repository to OCI mapping):
- `CommonTransportFormat`
- `DockerDaemon`
- `Empty`
- `OCIRegistry`

*Example:*
<pre>
$ ocm get componentversion ghcr.io/mandelsoft/kubelink
$ ocm get componentversion --repo OCIRegistry:ghcr.io mandelsoft/kubelink
</pre>


```
ocm components get [<options>] {<component-reference>} [flags]
```

### Options

```
  -h, --help               help for get
  -o, --output string      output mode (yaml, json, JSON)
  -r, --repo string        repository name or spec
  -s, --sort stringArray   sort fields
```

### SEE ALSO

* [ocm components](ocm_components.md)	 - 
