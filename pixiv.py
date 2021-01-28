import requests
import json
import re
import os


class Pixiv:
    def __init__(self):
        self.session = requests.Session()
        self.session.headers = {
            "accept":
            "text/html,application/xhtml+xml,application/xml;q=0.9,image/avif,image/webp,image/apng,*/*;q=0.8,application/signed-exchange;v=b3;q=0.9",
            "user-agent":
            "Mozilla/5.0 (X11; Linux x86_64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/88.0.4324.96 Safari/537.36",
            "cookie": ""
        }

    def get_ranklist(self):
        url = "https://www.pixiv.net/ranking.php?mode=daily&content=illust&p=1&format=json"
        res = self.session.get(url)
        data = json.loads(res.text)

        return data["contents"][:10]

    def get_images(self, illust_id):
        img_urls = []
        # pattern = r"/img/\d*.*/\d*_p0"
        pattern = r'(?<="original":").*?(?="\})'
        for i in illust_id:
            page_url = "https://www.pixiv.net/artworks/" + str(i)
            res = self.session.get(page_url)
            img_url = re.search(pattern, res.text, flags=0).group()
            img_urls.append(img_url)

        return img_urls

    def dl_images(self, urls, illust_id):
        # Get script folder path.
        file_path = os.path.dirname(os.path.abspath(__file__)) + "/rank_img"
        # Create folder.
        if not os.path.exists(file_path):
            os.mkdir(file_path)
        # Download images.
        index = 1
        for u in urls:
            f = self.session.get(u,
                                 headers={
                                     "referer":
                                     "https://www.pixiv.net/artworks/" +
                                     str(illust_id)
                                 })
            fm = re.search(r"[.](jpg|png|jpeg)$", u, flags=0).group()
            with open(file_path + "/" + str(index) + fm, "wb") as img:
                img.write(f.content)
            index += 1


if __name__ == "__main__":
    pixiv = Pixiv()
    ranklist = pixiv.get_ranklist()
    id_list = [l["illust_id"] for l in ranklist]
    img_urls = pixiv.get_images(id_list)
    pixiv.dl_images(img_urls, id_list)
