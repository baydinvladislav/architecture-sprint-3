package ru.yandex.practicum.smarthome.supplier;

import lombok.Getter;
import lombok.Setter;

@Setter
@Getter
public class TemperatureRequest {

    private String action;
    private double temperature;

    public TemperatureRequest(String action, double temperature) {
        this.action = action;
        this.temperature = temperature;
    }
}
