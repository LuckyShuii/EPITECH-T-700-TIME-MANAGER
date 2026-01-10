export const isTokenExpired = (token: string): boolean => {
  try {
    // Extraire la deuxième partie du JWT (le payload)
    const payloadBase64 = token.split('.')[1]
    
    // Décoder de base64
    const payloadJson = atob(payloadBase64)
    
    // Parser le JSON
    const payload = JSON.parse(payloadJson)
    
    // Vérifier si le token a expiré
    const expirationTime = payload.exp * 1000 // exp est en secondes, on convertit en ms
    return Date.now() >= expirationTime
  } catch (error) {
    // Si le token est invalide ou ne peut pas être décodé
    return true
  }
}