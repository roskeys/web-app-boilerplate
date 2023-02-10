from website.auth import auth
from dotenv import load_dotenv
from utils.constants import app
from flask_http_middleware import MiddlewareManager
from middlewares.auth_validator import AccessCheckerMiddleware

load_dotenv()
app.register_blueprint(auth, url_prefix="/auth")
app.wsgi_app = MiddlewareManager(app)
app.wsgi_app.add_middleware(AccessCheckerMiddleware)
