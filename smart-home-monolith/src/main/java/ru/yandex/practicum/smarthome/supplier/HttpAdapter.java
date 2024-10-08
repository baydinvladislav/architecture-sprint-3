package ru.yandex.practicum.smarthome.supplier;

import ru.yandex.practicum.smarthome.dto.HeatingSystemDto;

public interface HttpAdapter {
    HeatingSystemDto sendSetTemperatureRequest(Long moduleId, double temperature);
}
