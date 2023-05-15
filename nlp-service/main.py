from flask import Flask, request, jsonify
from nlp import translator, simplifier
from helpers import error_json

app = Flask(__name__)

@app.post('/process-text')
def process_request():
    process: function
    response_data: dict

    data = request.get_json()

    nlp_proccess = ''
    try:
        nlp_proccess = data['process']
    except:
        return error_json('Cannot retrieve parameter proccess from json')


    match nlp_proccess:
        case 'translate':
            process = translator.translate_to_el
        case 'simplify':
            process = simplifier.simplify
        case _:
            return error_json(f'Unrecognized proccess: {nlp_proccess}')

    try:
        desired_text = process(data['text'])

        response_data = {
            'error': False,
            'message': f'NLP proccess "{nlp_proccess}" was successful',
            'data': desired_text
        }
        return jsonify(response_data)
        
    except Exception as e:
        return error_json(str(e))
    


if __name__ == '__main__':
    app.run(debug=True, host='0.0.0.0', port=80)
