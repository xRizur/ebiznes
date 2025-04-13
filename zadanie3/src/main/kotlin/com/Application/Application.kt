package com.Application

import dev.kord.common.entity.Snowflake
import dev.kord.core.Kord
import dev.kord.core.behavior.channel.createMessage
import dev.kord.core.behavior.channel.TextChannelBehavior
import dev.kord.core.event.message.MessageCreateEvent
import dev.kord.core.on
import dev.kord.gateway.PrivilegedIntent
import dev.kord.common.entity.PresenceStatus
import dev.kord.gateway.Intent
import io.ktor.client.*
import io.ktor.client.engine.cio.*
import io.ktor.client.plugins.contentnegotiation.*
import io.ktor.client.request.*
import io.ktor.client.statement.*
import io.ktor.http.*
import io.ktor.serialization.kotlinx.json.*
import io.ktor.server.application.*
import io.ktor.server.engine.*
import io.ktor.server.netty.*
import io.ktor.server.plugins.contentnegotiation.*
import io.ktor.server.request.*
import io.ktor.server.response.*
import io.ktor.server.routing.*
import kotlinx.coroutines.*
import kotlinx.serialization.json.*
import kotlinx.coroutines.runBlocking

val KATEGORIE = listOf("Elektronika", "Książki", "Odzież")
val PRODUKTY = mapOf(
    "elektronika" to listOf("Smartfon", "Laptop", "Tablet", "Słuchawki", "Monitor"),
    "książki" to listOf("Powieść", "Biografia", "Fantastyka",),
    "odzież" to listOf("T-shirt", "Buty", "Czapka"),
)

fun main() {
    runBlocking {
        mainApp()
    }
}

suspend fun mainApp() = coroutineScope {
    val discordToken = System.getenv("DISCORD_TOKEN") ?: error("Brak zmiennej środowiskowej DISCORD_TOKEN")
    val discordChannelId = System.getenv("DISCORD_CHANNEL_ID") ?: error("Brak zmiennej środowiskowej DISCORD_CHANNEL_ID")
    val port = 8080

    @OptIn(PrivilegedIntent::class)
    val kord = Kord(discordToken)

    GlobalScope.launch {
        kord.on<MessageCreateEvent> {
            val content = message.content.trim()
            val parts = content.split(" ", limit = 3)
            val command = parts.firstOrNull()?.lowercase() ?: ""
            when (command) {
                "!help" -> {
                    val helpMessage = """
                    **Dostępne komendy:**
                    `!kategorie` - Pokaż wszystkie kategorie produktów
                    `!produkty [kategoria]` - Pokaż produkty w wybranej kategorii
                    `!help` - Wyświetl pomoc
                    """.trimIndent()
                    message.channel.createMessage(helpMessage)
                }
                
                "!kategorie" -> {
                    if (KATEGORIE.isEmpty()) {
                        message.channel.createMessage("Brak dostępnych kategorii.")
                    } else {
                        message.channel.createMessage("**Dostępne kategorie:** ${KATEGORIE.joinToString(", ")}")
                    }
                }
                
                "!produkty" -> {
                    if (parts.size < 2) {
                        message.channel.createMessage("Użycie: `!produkty [kategoria]`\nDostępne kategorie: ${KATEGORIE.joinToString(", ")}")
                        return@on
                    }
                    
                    val categoryName = parts[1].lowercase()
                    val products = PRODUKTY[categoryName] ?: emptyList()
                    
                    if (products.isEmpty()) {
                        message.channel.createMessage("Brak produktów dla kategorii: **$categoryName** lub nieprawidłowa kategoria.\nSprawdź `!kategorie`, aby zobaczyć dostępne kategorie.")
                    } else {
                        message.channel.createMessage("🛒 **Produkty w kategorii $categoryName:**\n${products.joinToString("\n") { "- $it" }}")
                    }
                }
            }
        }
        kord.login {
            @OptIn(PrivilegedIntent::class)
            intents += Intent.MessageContent
        }
    }

    embeddedServer(Netty, port = port) {
        install(io.ktor.server.plugins.contentnegotiation.ContentNegotiation) { 
            json() 
        }
        
        routing {
            post("/discord/send") {
                val params = call.receiveParameters()
                val message = params["message"]
                    ?: return@post call.respondText(
                        "Brak parametru 'message'", 
                        status = HttpStatusCode.BadRequest
                    )
                
                try {
                    val channel = kord.getChannel(Snowflake(discordChannelId)) as? TextChannelBehavior
                        ?: return@post call.respondText(
                            "Nie znaleziono kanalu o ID: $discordChannelId", 
                            status = HttpStatusCode.NotFound
                        )
                    
                    channel.createMessage(message)
                    call.respondText(
                        "Wiadomość została wysłana pomyślnie!", 
                        status = HttpStatusCode.OK
                    )
                } catch (e: Exception) {
                    call.respondText(
                        "Błąd podczas wysyłania wiadomości: ${e.message}", 
                        status = HttpStatusCode.InternalServerError
                    )
                }
            }
        }
    }.start(wait = true)
}
