# Bitbucket PR Reviewer Playbook

A step-by-step guide to reviewing a pull request on Bitbucket using the available MCP tools.

## Steps

1. **Find the pull request**
   Use `bitbucket_pr` with operation `listPullRequests` to list open pull requests for the repository. Note the `id` of the PR you want to review.

2. **Get PR details**
   Use `bitbucket_pr` with operation `getAPullRequest` and the PR `id` to read the title, description, source branch, destination branch, and author.

3. **Review the diff**
   Use `bitbucket_pr` with operation `listChangesInAPullRequest` to see which files changed. Use `getTheDiffStatForAPullRequest` for a summary of additions and deletions.

4. **Check commits**
   Use `bitbucket_pr` with operation `listCommitsOnAPullRequest` to review the individual commits in the PR.

5. **Read existing comments**
   Use `bitbucket_pr` with operation `listCommentsOnAPullRequest` to see what feedback has already been given.

6. **Leave a review comment**
   Use `bitbucket_pr` with operation `createACommentOnAPullRequest` to post your feedback. Provide the comment text in the `content.raw` body field.

7. **Approve or request changes**
   - To approve: use `bitbucket_pr` with operation `approveAPullRequest`.
   - To request changes: leave a comment explaining what needs to change.

## Tips

- Always read the PR description and existing comments before leaving feedback.
- Focus on correctness, readability, and potential edge cases.
- Be constructive and specific in your comments.
