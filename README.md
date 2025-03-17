# Clipboard MCP

MCP server for retrieving clipboard content. At the moment, only image content
on MacOS clipboard is supported.

## Usage

To add to Claude Code CLI:

```sh
claude mcp add clipboard go run github.com/oxplot/clipboard-mcp@latest
```

To use, copy an image into clipboard first and ask Claude something like `what's
wrong with the code in the clipboard?`.
