model "fast" {
  provider = "openai"
  id       = "gpt-4o-mini"
}

model "smart" {
  provider = "anthropic"
  id       = "claude-sonnet-5"
}

target "langgraph" {
  type   = "codegen"
  output = "./gen/langgraph"
}
