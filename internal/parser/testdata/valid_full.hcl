model "fast" {
  provider = "openai"
  id       = "gpt-4o-mini"

  params {
    temperature = 0.2
    max_tokens  = 4096
  }
}

model "smart" {
  provider = "anthropic"
  id       = "claude-sonnet-5"
}

target "langgraph" {
  type   = "codegen"
  output = "./gen/langgraph"
}

target "openai_assistants" {
  type = "platform"

  auth {
    api_key_env = "OPENAI_API_KEY"
  }
}
