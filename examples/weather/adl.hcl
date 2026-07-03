model "fast" {
  provider = "openai"
  id       = "gpt-4o-mini"

  params {
    temperature = 0.2
    max_tokens  = 4096
  }
}

# Codegen target -> exercises the module-walk skip of target output paths (#6)
target "langgraph" {
  type   = "codegen"
  output = "./gen/langgraph"
}

# Platform target with explicit auth (also valid without auth: ambient credentials)
target "openai_assistants" {
  type = "platform"

  auth {
    api_key_env = "OPENAI_API_KEY"
  }
}
