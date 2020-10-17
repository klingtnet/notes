let el = document.getElementById("note");
if (el) {
  el.addEventListener("keydown", function (e) {
    console.log(e);
    if (e.ctrlKey && e.key === "Enter") {
      this.form.submit();
    }
  });
}
