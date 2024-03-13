import os
import toml
import shutil

def main():
    parsed_toml = toml.load('../config.toml')
    sql_part = parsed_toml['sql']
    if sql_part is None:
        print('specify.py: missing part `sql` in config.toml')
        return -1
    tokens = ['username', 'password', 'database']
    replace_table = dict()
    for token in tokens:
        if sql_part[token] is None:
            print(f'specify.py: missing `{token}` in config.toml')
            return -1
        replace_table[token] = sql_part[token]
    
    DCYML_PATH = '../docker/docker-compose.yml'
    if os.path.exists(DCYML_PATH):
        os.remove(DCYML_PATH)

    with open(DCYML_PATH + '.template', 'r') as f:
        docker_compose_content = f.readlines()
    for (idx, line) in enumerate(docker_compose_content):
        for key, value in replace_table.items():
            line = line.replace('$' + key.upper(), value)
        docker_compose_content[idx] = line

    with open(DCYML_PATH, 'w+') as f:
        f.writelines(docker_compose_content)

if __name__ == '__main__':
    exit(main())

