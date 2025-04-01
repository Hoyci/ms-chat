import { format, isToday, isYesterday, differenceInHours } from "date-fns";
import { ptBR } from "date-fns/locale";

export function formatMessageDate(date: Date) {
  const messageDate = new Date(date);

  if (isToday(messageDate)) {
    return format(messageDate, "HH:mm");
  }

  if (isYesterday(messageDate)) {
    return "Ontem";
  }

  if (differenceInHours(new Date(), messageDate) < 168) {
    return format(messageDate, "EEEE", { locale: ptBR });
  }

  return format(messageDate, "dd/MM/yyyy");
}
