const $ = (id) => document.getElementById(id);

const urlInput = $("url");
const btnDownload = $("btn-download");
const errorSection = $("error-section");
const errorMsg = $("error-msg");
const btnRetry = $("btn-retry");

function getFormat() {
  const checked = document.querySelector('input[name="format"]:checked');
  return checked ? checked.value : "mp4";
}

function showError(msg) {
  errorMsg.textContent = msg;
  errorSection.classList.remove("hidden");
  btnDownload.disabled = false;
}

async function startDownload() {
  const url = urlInput.value.trim();
  if (!url) {
    showError("Informe a URL do YouTube.");
    return;
  }

  btnDownload.disabled = true;
  errorSection.classList.add("hidden");

  try {
    const res = await fetch("/api/downloads", {
      method: "POST",
      headers: { "Content-Type": "application/json" },
      body: JSON.stringify({ url, format: getFormat() }),
    });
    const data = await res.json();
    if (!res.ok) {
      showError(data.error || "Erro ao iniciar download.");
      return;
    }
    window.location.href = data.page_url || "/d/" + data.id;
  } catch (err) {
    showError("Falha de conexão com o servidor.");
  }
}

btnDownload.addEventListener("click", startDownload);
btnRetry.addEventListener("click", () => errorSection.classList.add("hidden"));
urlInput.addEventListener("keydown", (e) => {
  if (e.key === "Enter") startDownload();
});
