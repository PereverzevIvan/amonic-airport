import React from "react";
import "./Button.scss";

const ButtonColors = ["blue", "green", "red"];

function Button({
  children,
  color = "blue",
  onClick,
  type = "button",
  disabled = false,
}) {
  if (!ButtonColors.includes(color)) {
    console.log("Кнопка получила недопустимый цвет:", color);
    color = "blue";
  }

  return (
    <button
      type={type}
      onClick={onClick}
      className={`button button_color-${color}`}
      disabled={disabled}
    >
      {children}
    </button>
  );
}

export default Button;
