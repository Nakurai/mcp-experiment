# Session2: The demo server

As we saw in [session 1](/docs/session1.md), the whole purpose of an MCP server is to bridge an LLM with a tool that it can use.

Let's build a trivial version of such a tool. We will make a tool that can receive messages from an outside source and send messages to people.

By connecting the LLM to this tool, it can summarize our messages, suggest answers and send messages to people.

## Endpoints
The tool will be a demo server. Here are the endpoints:

- POST /api/login: initiate the Github Oauth web flow
- POST /messages: content,recipient,date / headers: access token
- GET /messages: date / headers: access token
- GET /search: content, recipient, date / headers: access token