from datetime import datetime
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

    def get_ranklist(self, restrict=False):
        url = "https://www.pixiv.net/ranking.php?mode={mode}&content=illust&p=1&format=json".format(
            mode="daily_r18" if restrict else "daily")
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

    def dl_images(self, dir_name, urls, illust_id):
        # Get date.
        now = datetime.now()
        date = now.strftime('%Y%m%d')
        # Get file path.
        file_path = os.path.join(os.path.dirname(os.path.abspath(__file__)),
                                 dir_name, date)
        # Create folder.
        if not os.path.exists(file_path):
            os.makedirs(file_path)
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
            with open(os.path.join(file_path, str(index) + fm), "wb") as img:
                img.write(f.content)
            index += 1


if __name__ == "__main__":
    pixiv = Pixiv()
    # Get list.
    ranklist, ranklist_r18 = pixiv.get_ranklist(), pixiv.get_ranklist(True)
    # Get id of artworks.
    id_list = [l["illust_id"] for l in ranklist]
    id_list_r18 = [l["illust_id"] for l in ranklist_r18]
    # Get url of images.
    img_urls = pixiv.get_images(id_list)
    img_urls_r18 = pixiv.get_images(id_list_r18)
    # Download
    pixiv.dl_images("rank_img", img_urls, id_list)
    pixiv.dl_images("rank_img_r18", img_urls_r18, id_list_r18)
