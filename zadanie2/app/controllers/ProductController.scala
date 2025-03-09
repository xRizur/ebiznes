package controllers

import javax.inject._
import play.api.mvc._
import play.api.libs.json._
import scala.collection.mutable.ListBuffer

case class Product(id: Long, name: String, price: Double)
object Product {
  implicit val productFormat: Format[Product] = Json.format[Product]
}

@Singleton
class ProductController @Inject()(cc: ControllerComponents) extends AbstractController(cc) {
  private val products = ListBuffer[Product](
    Product(1, "Product 1", 10.0),
    Product(2, "Product 2", 20.0)
  )

  def list: Action[AnyContent] = Action {
    Ok(Json.toJson(products))
  }

  def get(id: Long): Action[AnyContent] = Action {
    products.find(_.id == id).map { product =>
      Ok(Json.toJson(product))
    } getOrElse NotFound(Json.obj("error" -> "Product not found"))
  }

  def create: Action[JsValue] = Action(parse.json) { request =>
    request.body.validate[Product].map { product =>
      products += product
      Created(Json.toJson(product))
    } recoverTotal { _ =>
      BadRequest(Json.obj("error" -> "Invalid product format"))
    }
  }

  def update(id: Long): Action[JsValue] = Action(parse.json) { request =>
    request.body.validate[Product].map { updatedProduct =>
      products.indexWhere(_.id == id) match {
        case -1 => NotFound(Json.obj("error" -> "Product not found"))
        case idx =>
          products.update(idx, updatedProduct)
          Ok(Json.toJson(updatedProduct))
      }
    } recoverTotal { _ =>
      BadRequest(Json.obj("error" -> "Invalid product format"))
    }
  }

  def delete(id: Long): Action[AnyContent] = Action {
    val initialSize = products.size
    products.filterInPlace(_.id != id)
    if (products.size < initialSize)
      Ok(Json.obj("message" -> "Product deleted"))
    else
      NotFound(Json.obj("error" -> "Product not found"))
  }
}
