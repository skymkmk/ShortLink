import "./style.css";

const realURL: HTMLInputElement = document.getElementById(
  "realURL"
) as HTMLInputElement;
const button: HTMLButtonElement = document.getElementById(
  "button"
) as HTMLButtonElement;
const notice: HTMLDivElement = document.getElementById(
  "notice"
) as HTMLDivElement;
function handleClick() {
  let url = realURL.value;
  url = encodeURIComponent(url);
  const searchParam: URLSearchParams = new URLSearchParams({
    url: url,
  });
  fetch(
    `//${window.location.host}/api/v1/newShortLink?${searchParam.toString()}`,
    {
      method: "GET",
    }
  )
    .then((resp) => {
      resp
        .json()
        .then((data) => {
          if ((data?.status ?? -1) === 0) {
            notice.innerText = `短链为：${data?.realURL ?? ""}`;
          } else if ((data?.status ?? -1) === 3) {
            notice.innerText = `已生成过短链！${data?.realURL ?? ""}`;
          } else {
            notice.innerText = data?.error ?? "";
          }
          if ((data?.status ?? -1) === 0 || (data?.status ?? -1) === 3) {
            try {
              navigator.clipboard.writeText(data?.realURL ?? "");
              notice.innerText += "\n链接已复制至剪贴板！"
            } catch {}
          }
        })
        .catch((e) => {
          notice.innerText = e;
        });
    })
    .catch((e) => {
      notice.innerText = e;
    });
}
button.addEventListener("click", handleClick);
realURL.addEventListener("keydown", (e: KeyboardEvent) => {
  if (e.code === "Enter") {
    e.preventDefault();
    handleClick();
  }
});
