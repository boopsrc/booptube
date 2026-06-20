const $ = (id) => document.getElementById(id);

const formSection = $("form-section");
const progressSection = $("progress-section");
const readySection = $("ready-section");
const errorSection = $("error-section");

const urlInput = $("url");
const btnDownload = $("btn-download");
const statusText = $("status-text");
const progressBar = $("progress-bar");
const progressPct = $("progress-pct");
const logEl = $("log");
const filenameEl = $("filename");
const expiresEl = $("expires");
const downloadLink = $("download-link");
const errorMsg = $("error-msg");
const btnRetry = $("btn-retry");

let pollTimer = null;
let countdownTimer = null;
let currentJobId = null;

const statusLabels = {
  queued: "Na fila",
  downloading: "Baixando",
  ready: "Pronto",
  failed: "Falhou",
  expired: "Expirado",
};

function getFormat() {
  const checked = document.querySelector('input[name="format"]:checked');
  return checked ? checked.value : "mp4";
}

function showSection(section) {
  [formSection, progressSection, readySection, errorSection].forEach((el) => {
    el.classList.toggle("hidden", el !== section);
  });
}

function resetUI() {
  clearInterval(pollTimer);
  clearInterval(countdownTimer);
  pollTimer = null;
  countdownTimer = null;
  currentJobId = null;
  btnDownload.disabled = false;
  showSection(formSection);
}

function setProgress(pct) {
  const v = Math.min(100, Math.max(0, pct));
  progressBar.style.width = v + "%";
  progressPct.textContent = Math.round(v) + "%";
}

function setLog(lines) {
  if (!lines || !lines.length) {
    logEl.textContent = "";
    return;
  }
  logEl.textContent = lines.join("\n");
  logEl.scrollTop = logEl.scrollHeight;
}

async function startDownload() {
  const url = urlInput.value.trim();
  if (!url) {
    showError("Informe a URL do YouTube.");
    return;
  }

  btnDownload.disabled = true;
  errorSection.classList.add("hidden");
  showSection(progressSection);
  setProgress(0);
  setLog([]);
  statusText.textContent = statusLabels.queued;

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
    currentJobId = data.id;
    pollStatus();
    pollTimer = setInterval(pollStatus, 2000);
  } catch (err) {
    showError("Falha de conexão com o servidor.");
  }
}

async function pollStatus() {
  if (!currentJobId) return;

  try {
    const res = await fetch("/api/downloads/" + currentJobId);
    const job = await res.json();
    if (!res.ok) {
      showError(job.error || "Job não encontrado.");
      return;
    }

    statusText.textContent = statusLabels[job.status] || job.status;
    setProgress(job.progress || 0);
    setLog(job.log);

    if (job.status === "ready") {
      clearInterval(pollTimer);
      pollTimer = null;
      showReady(job);
    } else if (job.status === "failed") {
      clearInterval(pollTimer);
      pollTimer = null;
      showError(job.error || "Download falhou.");
    } else if (job.status === "expired") {
      clearInterval(pollTimer);
      pollTimer = null;
      showError("O arquivo expirou e foi removido do servidor.");
    }
  } catch (err) {
    // keep polling on transient errors
  }
}

function showReady(job) {
  filenameEl.textContent = job.filename || "arquivo";
  downloadLink.href = job.download_url;
  downloadLink.download = job.filename || "";

  if (job.expires_at) {
    startCountdown(new Date(job.expires_at));
  }

  showSection(readySection);
  btnDownload.disabled = false;
}

function startCountdown(expiresAt) {
  clearInterval(countdownTimer);

  function tick() {
    const remaining = expiresAt - Date.now();
    if (remaining <= 0) {
      expiresEl.textContent = "Arquivo expirado.";
      clearInterval(countdownTimer);
      return;
    }
    const min = Math.floor(remaining / 60000);
    const sec = Math.floor((remaining % 60000) / 1000);
    expiresEl.textContent =
      "Disponível por mais " +
      (min > 0 ? min + " min " : "") +
      sec +
      " s";
  }

  tick();
  countdownTimer = setInterval(tick, 1000);
}

function showError(msg) {
  errorMsg.textContent = msg;
  showSection(errorSection);
  btnDownload.disabled = false;
}

btnDownload.addEventListener("click", startDownload);
btnRetry.addEventListener("click", resetUI);

urlInput.addEventListener("keydown", (e) => {
  if (e.key === "Enter") startDownload();
});
