package ru.yandex.practicum.smarthome.supplier;

import lombok.RequiredArgsConstructor;
import org.slf4j.Logger;
import org.slf4j.LoggerFactory;
import org.springframework.http.HttpEntity;
import org.springframework.http.HttpHeaders;
import org.springframework.http.HttpMethod;
import org.springframework.http.MediaType;
import org.springframework.http.ResponseEntity;
import org.springframework.stereotype.Component;
import org.springframework.web.client.RestTemplate;
import ru.yandex.practicum.smarthome.dto.HeatingSystemDto;

@Component
@RequiredArgsConstructor
public class MicroserviceSupplierAdapter implements HttpAdapter {

    private final RestTemplate restTemplate;
    private static final String TELEMETRY_SERVICE_URL = "http://smarthome-url:8000";
    private static final Logger logger = LoggerFactory.getLogger(MicroserviceSupplierAdapter.class);

    public HeatingSystemDto SendSetTemperatureRequest(Long moduleId, double temperature) {
        TemperatureRequest temperatureRequest = new TemperatureRequest("set_temperature", temperature);
        String url = String.format("%s/modules/%d/house/%d", TELEMETRY_SERVICE_URL, moduleId, moduleId);
        logger.info("Sending request to URL: {}", url);
        HttpHeaders headers = new HttpHeaders();
        // TODO: здесь нужно проставить jwt токен
        headers.setContentType(MediaType.APPLICATION_JSON);
        HttpEntity<TemperatureRequest> requestEntity = new HttpEntity<>(temperatureRequest, headers);

        try {
            ResponseEntity<HeatingSystemDto> responseEntity = restTemplate.exchange(
                    url,
                    HttpMethod.POST,
                    requestEntity,
                    HeatingSystemDto.class
            );

            if (responseEntity.getStatusCode().is2xxSuccessful()) {
                logger.info("Successfully sent temperature request: {}", responseEntity.getBody());
                return responseEntity.getBody();
            } else {
                logger.error("Failed to send temperature request: {}", responseEntity.getStatusCode());
                throw new RuntimeException("Failed to send temperature request. " +
                        "HTTP Status Code: " + responseEntity.getStatusCode());
            }
        } catch (Exception e) {
            logger.error("Error occurred while sending temperature request: {}", e.getMessage(), e);
            throw new RuntimeException("Error occurred while sending temperature request: " + e.getMessage(), e);
        }
    }
}
