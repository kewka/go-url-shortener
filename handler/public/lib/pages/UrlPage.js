(function () {
  const inputUrl = document.querySelector(".UrlPage-url-input");
  const btnUrlCopy = document.querySelector(".UrlPage-url-copy");
  btnUrlCopy.addEventListener("click", function () {
    inputUrl.select();
    document.execCommand("copy");
    btnUrlCopy.textContent = "Copied!";
  });
})();
