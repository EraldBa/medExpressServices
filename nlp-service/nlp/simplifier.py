from transformers import AutoTokenizer, AutoModelForCausalLM
import torch

TOKENIZER = AutoTokenizer.from_pretrained("philippelaban/keep_it_simple")
KIS_MODEL = AutoModelForCausalLM.from_pretrained("philippelaban/keep_it_simple")
MAX_SIMPLIFY_LEN = 2500

def simplify(text: str) -> str:
    if len(text) > MAX_SIMPLIFY_LEN:
        raise Exception('Text too big for simplification. Proccess aborted')

    start_id = TOKENIZER.bos_token_id

    tokenized_paragraph = [(TOKENIZER.encode(text=text) + [start_id])]

    input_ids = torch.LongTensor(tokenized_paragraph)

    output_ids = KIS_MODEL.generate(input_ids, max_length=len(text), num_beams=1, do_sample=True, num_return_sequences=4)

    output_ids = output_ids[:, input_ids.shape[1]:]

    output = TOKENIZER.batch_decode(output_ids)

    output = [out.replace(TOKENIZER.eos_token, "") for out in output]

    max_out = ''
    for out in output:
        max_out = max(max_out, out, key=len)

    return max_out