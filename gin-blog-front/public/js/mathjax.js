// 设置MathJax的配置
window.MathJax = {
    // 行内公式使用的开始和结束标记。
    tex: {
      // 行内公式选择符
      inlineMath: [
        ['$', '$'],
        ['\\(', '\\)'],
      ],
      // 段内公式选择符
      displayMath: [
        ['$$', '$$'],
        ['\\[', '\\]'],
      ],
    },
    // startup对象中的ready函数会在MathJax加载完成后被调用
    startup: {
      ready() {
        // 当MathJax加载完成后，调用默认的ready函数
        MathJax.startup.defaultReady()
      },
    },
  }
  
// 用来配置MathJax的。MathJax是一个JavaScript库，用于在网页上显示数学公式。
// 它支持LaTeX、MathML和AsciiMath等多种数学标记语言。