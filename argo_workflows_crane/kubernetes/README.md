If your using Helm to deploy your Argo Workflows Workflows your going to need to escape the variables.

```
"{{ `{{ workflow.parameters.build_path }}` }}"
```
