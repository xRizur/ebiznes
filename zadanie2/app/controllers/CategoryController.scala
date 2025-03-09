package controllers

import javax.inject._
import play.api.mvc._
import play.api.libs.json._
import scala.collection.mutable.ListBuffer

case class Category(id: Long, name: String)
object Category {
  implicit val categoryFormat: Format[Category] = Json.format[Category]
}

@Singleton
class CategoryController @Inject()(cc: ControllerComponents) extends AbstractController(cc) {
  private val categories = ListBuffer[Category](
    Category(1, "Category 1"),
    Category(2, "Category 2")
  )

  def list: Action[AnyContent] = Action {
    Ok(Json.toJson(categories))
  }

  def get(id: Long): Action[AnyContent] = Action {
    categories.find(_.id == id).map { category =>
      Ok(Json.toJson(category))
    } getOrElse NotFound(Json.obj("error" -> "Category not found"))
  }

  def create: Action[JsValue] = Action(parse.json) { request =>
    request.body.validate[Category].map { category =>
      categories += category
      Created(Json.toJson(category))
    } recoverTotal { _ =>
      BadRequest(Json.obj("error" -> "Invalid category format"))
    }
  }

  def update(id: Long): Action[JsValue] = Action(parse.json) { request =>
    request.body.validate[Category].map { updatedCategory =>
      categories.indexWhere(_.id == id) match {
        case -1 => NotFound(Json.obj("error" -> "Category not found"))
        case idx =>
          categories.update(idx, updatedCategory)
          Ok(Json.toJson(updatedCategory))
      }
    } recoverTotal { _ =>
      BadRequest(Json.obj("error" -> "Invalid category format"))
    }
  }

  def delete(id: Long): Action[AnyContent] = Action {
    val initialSize = categories.size
    categories.filterInPlace(_.id != id)
    if (categories.size < initialSize)
      Ok(Json.obj("message" -> "Category deleted"))
    else
      NotFound(Json.obj("error" -> "Category not found"))
  }
}
