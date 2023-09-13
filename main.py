import json
import datetime as dt

def string_santizer(data:dict):
    return data['S'].strip()

def bool_sanitizer(bool_data):
    if bool_data in ['1','t','T', 'TRUE', 'true', 'True']:
        return True
    elif bool_data in ['0','f','F', 'FALSE', 'false', 'False']:
        return False
    else:
        return None

def number_sanitizer(data: dict):
    try:
        return float(data['N'])
    except ValueError as e:
        return None


def sanitify(data: dict):
    if data:
        if 'S' in data:
            try:
                timestamp = dt.datetime.strptime(data['S'], '%Y-%m-%dT%H:%M:%SZ')
                return int(timestamp.timestamp()), True
            except ValueError as e:
                return string_santizer(data), True
        elif 'BOOL' in data:
            return bool_sanitizer(data['BOOL'].strip()), True
        elif 'NULL' in data:
            response = bool_sanitizer(data['NULL'].strip())
            if response:
                return None, False
            else:
                return None, True
        elif 'N' in data:
            return number_sanitizer(data), True
        elif 'L' in data and type(data['L']) is list:
            result = []
            for x in data['L']:
                list_sanitized_data, override = sanitify(x)
                if override and list_sanitized_data is not None and list_sanitized_data != '':
                    result.append(list_sanitized_data)
                elif not override:
                    result.append(list_sanitized_data)
            return result, True

        elif 'M' in data:
            inlinemap = {}
            for k, v in data['M'].items():
                inlinesanitifydata, override = sanitify(v)
                if k:
                    if override and inlinesanitifydata is not None:
                        inlinemap[k] = inlinesanitifydata
                    elif not override:
                        inlinemap[k] = inlinesanitifydata
            return inlinemap, True
        else:
            return None, True


def sanitize():
    response = {}
    with open('input.json','r') as r:
        data = json.load(r)
        for k, v in data.items():
            if k:
                sanitized_data, override = sanitify(v)
                if override and sanitized_data is not None:
                    response[k] = sanitized_data
                elif not override:
                    response[k] = sanitized_data
        print(json.dumps(response))





if __name__=="__main__":
    sanitize()