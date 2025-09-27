function adjustIframeHeight() {
  try {
    const iframe = window.frameElement
    if (!iframe) return
    const height = Math.max(
      document.body.scrollHeight,
      document.body.offsetHeight,
      document.documentElement.clientHeight,
      document.documentElement.scrollHeight,
      document.documentElement.offsetHeight,
    )
    iframe.style.height = height + 'px'
  } catch (e) {
    console.error('调整 iframe 高度失败:', e)
  }
}
window.addEventListener('load', adjustIframeHeight)

const observer = new MutationObserver(adjustIframeHeight)
observer.observe(document.body, {
  attributes: true,
  childList: true,
  subtree: true,
  characterData: true,
})
window.addEventListener('resize', adjustIframeHeight)
