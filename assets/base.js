let el = document.getElementById("note");
if (el) {
  el.addEventListener("keydown", function (e) {
    if (e.ctrlKey && e.key === "Enter") {
      this.form.submit();
    }
  });
}
