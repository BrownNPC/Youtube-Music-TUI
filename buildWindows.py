import hashlib
import os
import time
from bs4 import BeautifulSoup
import requests
def download_file(url, file_name):
    """
    Downloads a file from a given URL and saves it locally.
    """
    response = requests.get(url, stream=True)
    with open(file_name, 'wb') as file:
        for chunk in response.iter_content(chunk_size=1024):
            if chunk:
                file.write(chunk)

def get_md5(file_name):
    """
    Computes the MD5 checksum of a file.
    """
    md5_hash = hashlib.md5()
    with open(file_name, "rb") as f:
        for byte_block in iter(lambda: f.read(4096), b""):
            md5_hash.update(byte_block)
    return md5_hash.hexdigest()

scraped_links = []
def get_package(url):
    """
    Scrapes the given URL, downloads the file, and saves the MD5 checksum.
    Also downloads dependencies if any are found.

    Args:
    url (str): The URL of the page to scrape.

    Returns:
    str: The 'href' of the first "File" link found, or an empty string if none found.
    """
    print(f"Scraping package page: {url}")
    
    # Send a GET request to the page
    response = requests.get(url)

    # Check if the request was successful
    if response.status_code == 200:
        # Parse the page content using BeautifulSoup
        soup = BeautifulSoup(response.content, 'html.parser')
        # Find all 'dt' elements with the text "File:"
        dt_elements = soup.find_all('dt', string="File:")
        file_link = ''
        # Loop through each 'dt' element
        for dt in dt_elements:
            # Find the next sibling 'dd' tag that contains the link
            dd = dt.find_next_sibling('dd')
            if dd:
                # Find the 'a' tag inside the 'dd' tag and extract the 'href'
                link = dd.find('a')
                if link and 'href' in link.attrs:
                    file_link = link['href']
                    
                    # Look for the md5 checksum two elements after the file link (below SHA256)
                    sha256_dt = dt.find_next_siblings('dt', string="SHA256:")
                    if sha256_dt:
                        md5_dd = sha256_dt[0].find_next_sibling('dd')
                        if md5_dd:
                            md5_checksum = md5_dd.get_text(strip=True)
                    break  # Break after finding the first valid link

        if file_link:
            file_name = file_link.split('/')[-1]

            # Check if file already exists
            if os.path.exists(file_name):
                print(f"{file_name} already exists, skipping download.")
            else:
                print(f"Downloading {file_name}...")
                download_file(file_link, file_name)
                print(f"{file_name} downloaded.")

            # Save MD5 checksum to a file
            md5_from_file = get_md5(file_name)
            with open(f"{file_name}.md5", "w") as md5_file:
                md5_file.write(f"MD5 ({file_name}) = {md5_from_file}\n")
                print(f"MD5 checksum saved: {md5_from_file}")
        return file_link
    else:
        print(f"Failed to retrieve the page. Status code: {response.status_code}")
        return ''


import requests
from threading import Thread

downloaded_mpv = False
def download_latest_asset(owner="shinchiro", repo="mpv-winbuild-cmake"):
    global downloaded_mpv
    # GitHub API URL for the releases
    url = f"https://api.github.com/repos/{owner}/{repo}/releases/latest"
    
    # Get the latest release information
    response = requests.get(url)
    if response.status_code != 200:
        print(f"Failed to fetch releases: {response.status_code}")
        return

    release_data = response.json()

    # Automatically grab the asset name
    asset_name = next((asset['name'] for asset in release_data['assets']
                       if asset['name'].startswith("mpv-dev-x86_64-gcc")), None)
    
    if asset_name:
        # Correctly get the asset download URL
        asset_url = next(asset['url'] for asset in release_data['assets'] if asset['name'] == asset_name)
        
        # Download the asset
        asset_response = requests.get(asset_url, headers={"Accept": "application/octet-stream"})
        if asset_response.status_code == 200:
            with open(asset_name, 'wb') as f:
                f.write(asset_response.content)
            print(f"Downloaded: {asset_name}")
            downloaded_mpv = True
        else:
            print(f"Failed to download {asset_name}: {asset_response.status_code}")
    else:
        print("No asset matching the criteria found in the latest release.")

# Call the function


def go_to_path(path):
    if not os.path.exists(path):
        os.mkdir(path)
    os.chdir(path)
# if we are not in build folder, go there
go_to_path("build")

Thread(target=download_latest_asset).start()
get_package("https://packages.msys2.org/packages/mingw-w64-x86_64-mpv")
while not downloaded_mpv:
    time.sleep(1)
extracted_mpv = False
mpv_7z_path = ''
files = os.listdir()
for file in files:
    if file.endswith(".pkg.tar.zst"):
        print(f"\n\nextracting {file}\n--------------------------------------")
        os.system(f"tar -xvf {file}")
    elif file.endswith(".7z"):
        print(f"trying to extract mpv 7z\nthis will fail if you dont have 7z installed")
        try:
            os.system(f"7z x {file}")
            extracted_mpv = True
        except Exception as e:
            print(f"failed to extract {file}")
            print("either install p7zip or extract manually")
            extracted_mpv = False
            mpv_7z_path = file

os.system("clear")
mingw64_path = os.path.join(os.getcwd(),"mingw64")
PKG_CONFIG_PATH = os.path.join(mingw64_path,"lib","pkgconfig")
CGO_LDFLAGS = f"-L{mingw64_path}/lib"
CGO_CFLAGS = f"-I{mingw64_path}/include"
build_command = f"""
go mod tidy
export PKG_CONFIG_PATH={PKG_CONFIG_PATH}
export CGO_CFLAGS={CGO_CFLAGS} 
export CGO_LDFLAGS={CGO_LDFLAGS}
export CGO_ENABLED=1
export GOOS=windows
export CC=x86_64-w64-mingw32-gcc 
export CXX=x86_64-w64-mingw32-g++ 
go build -o ytt.exe ../ytt 
"""
go_to_path("../dist-windows")
print("Executable path:",os.path.join(os.getcwd(), "ytt.exe"))
os.system(build_command)
dist_path = os.getcwd()

go_to_path(mingw64_path)
if extracted_mpv:
    for file in os.listdir("../"):
        if file.endswith(".dll"):
            os.system(f"cp ../{file} {dist_path}")
else:
    print("mpv not extracted, please extract manually")
    print(os.path.join(mingw64_path, mpv_7z_path))
