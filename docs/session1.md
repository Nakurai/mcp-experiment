# MCP Server Basics

## Why an MCP server?

Okay so from what I understand, an MCP server is a server that will gives LLMs the ability to take some actions.

For example, right now, if one asks their favorite chat bot to order some food at a specific restaurant, it's not possible.

So here is the idea: the restaurant has a server that the bot can "talk" to, and this server will take care of ordering the food.

See, LLMs are great at a very limited set of tasks. For everything else, regular old business logic is still more efficient, simple, and also already done. So MCP is trying to create a bridge between LLMs and existing actions.

## What an MCP server?

MCP stands for Model Context Protocol. It's a big spec file that defines how the bot can talk to the server. You can find it on the [modelcontextprotocol](https://modelcontextprotocol.info) website

Well, almost. If I understand correctly, the bot will talk to an MCP "client", which in turn will be able to contact the MCP server.

The example I have seen around is Claude desktop. This would be the bot. Using the interface, you can install plugins (the clients) that will contact the MCP servers.

One last nuance is that a "server" in this context does not necessarily mean a machine running on the internet. It could be a software running on your computer that the plugin assumes is there and talk to as well.

## Who an MCP server?

MCP servers are popping up a bit everywhere, since it's very convenient to be able to ask LLMs to do task in natural language.

Naturally, lists/repositories are made. For examples:
- https://mcpservers.org/
- https://github.com/modelcontextprotocol/servers
- https://mcprepository.net/

In the future, I am guessing there will probably be an equivalent of Google search but for LLMs! the LLMs will be able to talk to an MCP server that will know what other MCP servers can do because they asked it to reference them. That could be a lucrative business.

## How an MCP server?

It seems like a big unknown here lies in the authentication. At least when it comes to MCP servers available via HTTP. The MCP spec says that the servers must implement OAuth 2.0 and the clients...well, nothing much there. At least not in terms of communication between the host and the client.

The question is this: where do you login? Or type your credit card information? In the example I gave earlier, I am guessing this will be a necessary feature of the plugins you would install in softwares like Claude Desktop (which also means that those softwares will need to give the plugins the ability to do so). They will need a secure way to receive sensitive data and make the LLMs able to use the resulting authentication mechanism (tokens etc.).

On this note, [this PR](https://github.com/modelcontextprotocol/modelcontextprotocol/pull/284) on the modelcontextprotocol github repo seems to indicates that the MCP server itself do not have to implement OAuth anymore, it can just be used as a resource server.

Alright, let's see what a basic MCP server looks like in [session 2](/docs/session2.md).


