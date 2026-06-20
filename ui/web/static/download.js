const $ = (id) => document.getElementById(id);

const loadingSection = $("loading-section");
const progressSection = $("progress-section");
const readySection = $("ready-section");
const errorSection = $("error-section");

const statusText = $("status-text");
const progressBar = $("progress-bar");
const progressPct = $("progress-pct");
const logEl = $("log");
const filenameEl = $("filename");
const expiresEl = $("expires");
const downloadLink = $("download-link");
const errorMsg = $("error-msg");

let pollTimer = null;
let countdownTimer = null;

const statusLabels = {
  queued: "Na fila",
  downloading: "Baixando",
  ready: "Pronto",
  failed: "Falhou",
  expired: "Expirado",
};

function jobIDFromPath() {
  const parts = window.location.pathname.split("/").filter(Boolean);
  if (parts[0] === "d" && parts[1]) {
    return parts[1];
  }
  return null;
}

const jobId = jobIDFromPath();

function showSection(section) {
  [loadingSection, progressSection, readySection, errorSection].forEach((el) => {
    el.classList.toggle("hidden", el !== section);
  });
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

function showError(msg) {
  clearInterval(pollTimer);
  clearInterval(countdownTimer);
  errorMsg.textContent = msg;
  showSection(errorSection);
}

async function pollStatus() {
  if (!jobId) return;

  try {
    const res = await fetch("/api/downloads/" + jobId);
    const job = await res.json();
    if (!res.ok) {
      showError(job.error || "Download não encontrado.");
      return;
    }

    if (job.status === "queued" || job.status === "downloading") {
      showSection(progressSection);
      statusText.textContent = statusLabels[job.status] || job.status;
      setProgress(job.progress || 0);
      setLog(job.log);
      return;
    }

    if (job.status === "ready") {
      clearInterval(pollTimer);
      pollTimer = null;
      showReady(job);
      return;
    }

    if (job.status === "failed") {
      showError(job.error || "Download falhou.");
      return;
    }

    if (job.status === "expired") {
      showError("O arquivo expirou e foi removido do servidor.");
    }
  } catch (err) {
    // keep polling on transient errors
  }
}

function showReady(job) {
  filenameEl.textContent = job.filename || "arquivo";
  downloadLink.href = job.download_url || "/d/" + jobId + "/file";
  downloadLink.download = job.filename || "";

  if (job.expires_at) {
    startCountdown(new Date(job.expires_at));
  }

  showSection(readySection);
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

if (!jobId) {
  showError("Link de download inválido.");
} else {
  showSection(loadingSection);
  pollStatus();
  pollTimer = setInterval(pollStatus, 1000);
}
