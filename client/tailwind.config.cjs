// tailwind.config.cjs
module.exports = {
  theme: {
    extend: {
      colors: {
        // Cores principais do WhatsApp
        primary: {
          50: "#dcf8c6", // Verde claro (para fundos de mensagens)
          100: "#202C33", // Verde principal do WhatsApp
          200: "#128C7E", // Verde escuro (hover/estações)
          300: "#075E54", // Verde muito escuro (cabeçalhos)
        },
        secondary: {
          50: "#34B7F1", // Azul claro (links/ícones)
          100: "#008069", // Azul-verde (destaques)
        },
        // Neutros e base
        neutral: {
          50: "#374248",
          100: "#AEBAC1", // Branco puro
          200: "#f8f9fa", // Fundo claro
          300: "#f0f2f5", // Bordas claras
          400: "#e9edef", // Fundo de input
          500: "#667781", // Texto secundário
          600: "#3b4a54", // Texto principal
          700: "#111b21", // Cabeçalhos escuros
        },
        background: "#262524",
        // Estados e extras
        success: "#25D366", // Verde de confirmação
        warning: "#FFD700", // Amarelo para avisos
        error: "#FF3B30", // Vermelho para erros
        info: "#34B7F1", // Azul para informações
      },
    },
  },
};
