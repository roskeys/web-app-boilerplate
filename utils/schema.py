email_regex = r"\b[A-Za-z0-9._%+-]+@[A-Za-z0-9.-]+\.[A-Z|a-z]{2,7}\b"

signup_schema = {
    "type": "object",
    "properties": {
        "email": {"type": "string"},
        "username": {"type": "string"},
        "password": {"type": "string"},
    },
    "required": ["email", "username", "password"],
}

login_schema = {
    "type": "object",
    "properties": {
        "email": {"type": "string"},
        "password": {"type": "string"},
    },
    "required": ["email", "password"],
}
