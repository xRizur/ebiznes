package controllers

import javax.inject._
import play.api.mvc._
import play.api.libs.json._
import scala.collection.mutable.ListBuffer

case class CartItem(id: Long, productId: Long, quantity: Int)
object CartItem {
  implicit val cartItemFormat: Format[CartItem] = Json.format[CartItem]
}

@Singleton
class CartController @Inject()(cc: ControllerComponents) extends AbstractController(cc) {
  private val cartItems = ListBuffer[CartItem](
    CartItem(1, 1, 2),
    CartItem(2, 2, 1)
  )

  def list: Action[AnyContent] = Action {
    Ok(Json.toJson(cartItems))
  }

  def get(id: Long): Action[AnyContent] = Action {
    cartItems.find(_.id == id).map { item =>
      Ok(Json.toJson(item))
    } getOrElse NotFound(Json.obj("error" -> "Cart item not found"))
  }

  def create: Action[JsValue] = Action(parse.json) { request =>
    request.body.validate[CartItem].map { item =>
      cartItems += item
      Created(Json.toJson(item))
    } recoverTotal { _ =>
      BadRequest(Json.obj("error" -> "Invalid cart item format"))
    }
  }

  def update(id: Long): Action[JsValue] = Action(parse.json) { request =>
    request.body.validate[CartItem].map { updatedItem =>
      cartItems.indexWhere(_.id == id) match {
        case -1 => NotFound(Json.obj("error" -> "Cart item not found"))
        case idx =>
          cartItems.update(idx, updatedItem)
          Ok(Json.toJson(updatedItem))
      }
    } recoverTotal { _ =>
      BadRequest(Json.obj("error" -> "Invalid cart item format"))
    }
  }

  def delete(id: Long): Action[AnyContent] = Action {
    val initialSize = cartItems.size
    cartItems.filterInPlace(_.id != id)
    if (cartItems.size < initialSize)
      Ok(Json.obj("message" -> "Cart item deleted"))
    else
      NotFound(Json.obj("error" -> "Cart item not found"))
  }
}
