import urllib.request
import os
from urllib.parse import urljoin


def download_file(url, filename):
    file_url = urljoin(url, filename)
    dataPath = "data"
    os.makedirs(dataPath, exist_ok=True)
    dest = os.path.join(dataPath, filename)
    try:
        urllib.request.urlretrieve(file_url, dest)
        print(f"File downloaded successfully: {filename}")
        return dest
    except Exception as e:
        print(f"Error downloading file: {e}")
        return None


def read_file(filename):
    dataPath = "data/"
    try:
        with open(filename, 'r') as file:
            data = file.read()
            print(f"File read successfully: {filename}")
            return data 
    except Exception as e:
        print(f"Error reading file: {e}")
        return None