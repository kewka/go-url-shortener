(function () {
  const themes = {
    light: "light",
    dark: "dark",
  };

  function getTheme() {
    return themes[localStorage.getItem("theme")] || themes.light;
  }

  function setTheme(theme) {
    localStorage.setItem("theme", theme);
    document.documentElement.setAttribute("data-theme", theme);
  }

  function toggleTheme() {
    setTheme(getTheme() === themes.light ? themes.dark : themes.light);
  }

  setTheme(getTheme());

  document.addEventListener("DOMContentLoaded", function () {
    document
      .querySelector(".Layout-toggle-theme")
      .addEventListener("click", toggleTheme);
  });
})();
