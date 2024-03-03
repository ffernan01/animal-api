# Repository config

With the aim of achieving more reliable code and to protect the branch flow, some configs where added to the repository. This is what I've done:

* Settings->Rules/Rulesets
* Target branches: `Default (main)`, `releases/**/*`
* Restrict deletions :white_check_mark:
* Require pull request before merging :white_check_mark:
    * Dismiss stale pull request approvals when new commits are pushed :white_check_mark:
* Require status check to pass :white_check_mark:
    * Add checks->`test-coverage`
* Block force pushes :white_check_mark: