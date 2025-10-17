const minutesToHours = (minutes: number): number => {
  return Math.round((minutes / 60) * 100) / 100 // Arrondi à 2 décimales
}