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

val CATEGORIES = listOf("Electronics", "Books", "Clothing")
val PRODUCTS = mapOf(
    "electronics" to listOf("Smartphone", "Laptop", "Tablet", "Headphones", "Monitor"),
    "books" to listOf("Novel", "Biography", "Fantasy",),
    "clothing" to listOf("T-shirt", "Shoes", "Cap"),
)

fun main() {
    runBlocking {
        mainApp()
    }
}

suspend fun mainApp() = coroutineScope {
    val discordToken = System.getenv("DISCORD_TOKEN") ?: error("Missing DISCORD_TOKEN environment variable")
    val discordChannelId = System.getenv("DISCORD_CHANNEL_ID") ?: error("Missing DISCORD_CHANNEL_ID environment variable")
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
                    **Available commands:**
                    `!categories` - Show all product categories
                    `!products [category]` - Show products in the selected category
                    `!help` - Display this help
                    """.trimIndent()
                    message.channel.createMessage(helpMessage)
                }
                
                "!categories" -> {
                    if (CATEGORIES.isEmpty()) {
                        message.channel.createMessage("No available categories.")
                    } else {
                        message.channel.createMessage("**Available categories:** ${CATEGORIES.joinToString(", ")}")
                    }
                }
                
                "!products" -> {
                    if (parts.size < 2) {
                        message.channel.createMessage("Usage: `!products [category]`\nAvailable categories: ${CATEGORIES.joinToString(", ")}")
                        return@on
                    }
                    
                    val categoryName = parts[1].lowercase()
                    val products = PRODUCTS[categoryName] ?: emptyList()
                    
                    if (products.isEmpty()) {
                        message.channel.createMessage("No products found for category: **$categoryName** or invalid category.\nCheck `!categories` to see available categories.")
                    } else {
                        message.channel.createMessage("🛒 **Products in category $categoryName:**\n${products.joinToString("\n") { "- $it" }}")
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
                        "Missing 'message' parameter", 
                        status = HttpStatusCode.BadRequest
                    )
                
                try {
                    val channel = kord.getChannel(Snowflake(discordChannelId)) as? TextChannelBehavior
                        ?: return@post call.respondText(
                            "Channel not found with ID: $discordChannelId", 
                            status = HttpStatusCode.NotFound
                        )
                    
                    channel.createMessage(message)
                    call.respondText(
                        "Message sent successfully!", 
                        status = HttpStatusCode.OK
                    )
                } catch (e: Exception) {
                    call.respondText(
                        "Error sending message: ${e.message}", 
                        status = HttpStatusCode.InternalServerError
                    )
                }
            }
        }
    }.start(wait = true)
}
