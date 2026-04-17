from openai import OpenAI
from typing import Optional


class LLMClient:
    def __init__(self, api_key: str, model: str = "gpt-4o", base_url: Optional[str] = None):
        if base_url:
            self.client = OpenAI(api_key=api_key, base_url=base_url)
        else:
            self.client = OpenAI(api_key=api_key)
        self.model = model

    def generate_test(self, code: str, language: str, category: str) -> str:
        prompt = self._build_prompt(code, language, category)
        response = self.client.chat.completions.create(
            model=self.model,
            messages=[
                {"role": "system", "content": "You are an expert software testing engineer. Generate comprehensive unit tests for the given code."},
                {"role": "user", "content": prompt}
            ],
            temperature=0.2,
            max_tokens=2000
        )
        return response.choices[0].message.content.strip()

    def _build_prompt(self, code: str, language: str, category: str) -> str:
        lang_map = {
            "python": "Python (use unittest or pytest)",
            "java": "Java (use JUnit)",
            "cpp": "C++ (use Google Test)",
            "go": "Go (use testing package)",
            "javascript": "JavaScript (use Jest)"
        }
        framework = lang_map.get(language, language)
        
        return f"""You are tasked with generating comprehensive unit tests for the following {language} code.

Category: {category}
Testing Framework: {framework}

SOURCE CODE:
```{language}
{code}
```

Requirements:
1. Generate comprehensive unit tests that cover:
   - Normal/happy path cases
   - Edge cases and boundary conditions
   - Error handling (if applicable)
   - Different input types and values

2. Follow best practices for {framework}:
   - Clear test function/method names
   - Proper assertions
   - Test isolation

3. The tests should be complete and runnable.

Please output ONLY the test code, wrapped in a code block with language tag."""


def create_openai_client(api_key: str, model: str = "gpt-4o") -> LLMClient:
    return LLMClient(api_key=api_key, model=model)


def create_deepseek_client(api_key: str, model: str = "deepseek-chat") -> LLMClient:
    return LLMClient(api_key=api_key, model=model, base_url="https://api.deepseek.com")


def create_ollama_client(model: str = "qwen2.5-coder") -> LLMClient:
    return LLMClient(api_key="ollama", model=model, base_url="http://localhost:11434/v1")
