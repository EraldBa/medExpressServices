from flask import jsonify, Response

def error_json(message: str) -> Response:
    err =  {
        'error': True,
        'message': message
    }

    return jsonify(err)