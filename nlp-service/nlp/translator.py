from transformers import pipeline

MAX_TEXT_LENGHT = 400
TRANSLATION_PIPELINE = pipeline("translation", model="Helsinki-NLP/opus-mt-en-el")

def translate_to_el(text: str) -> str:
    translated_text = ''

    try:
        text = str(text)
        text_slice = ''

        for i in range(0, len(text), MAX_TEXT_LENGHT):
            if len(text[i:]) < MAX_TEXT_LENGHT:
                text_slice = text[i:]
            else:
                text_slice = text[i:i+MAX_TEXT_LENGHT]
            
            translated_text += TRANSLATION_PIPELINE(text_slice, max_length=MAX_TEXT_LENGHT)[0]['translation_text']   
            
    except Exception as e:
        raise Exception(f'Could not translate provided text with error {str(e)}')
    
    return translated_text
    
