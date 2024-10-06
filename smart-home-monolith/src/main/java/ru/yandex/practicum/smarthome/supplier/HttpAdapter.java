package ru.yandex.practicum.smarthome.supplier;

import ru.yandex.practicum.smarthome.dto.HeatingSystemDto;

public interface HttpAdapter {
    HeatingSystemDto SendSetTemperatureRequest(Long moduleId, double temperature);
}
