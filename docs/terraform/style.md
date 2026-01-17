# Donald's terraform style guide

## PR tools

These are tools that can help with things like auto plan in PR, generate readable
diffs, change comparisons.

Something I've been trying to create is an algorithm to score changes. For
example: If I create a new tf module and run it, the algo reads the plan output
and git diff. From there it creates a score based on how many resources are
added (net new), changed, or deleted compared to the overall resource count. The
downside to this is building a weighted scorecard on resource types and context.

Weights:

1. Impact to the current run - Term: Local Blast Radius.
    - Example: if the IAM role breaks in the apply, the instance that uses that role fails as well.

2. Importance to the user - Term: Global Blast Radius.  
    - Example: This is a core piece of infra so its a 10/10 vs a new instance that's a 3/10 or whatever.

3. Policy checks with Open Policy Agent/Conftest. Term: Policy Checks.
    - Each check can be labeled with a score to generate an overall score.
    - A lot of policy checks are typically pass/fail. This might not be directly
      applicable but might be able to write rules in a way that can provide some
      usefulness.
