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
    private static final String MICROSERVICE_SERVICE_URL = "http://smarthome-url:8000";
    private static final Logger logger = LoggerFactory.getLogger(MicroserviceSupplierAdapter.class);

    public HeatingSystemDto getHeatingSystem(Long id) {
        String url = String.format("%s/modules/%d", MICROSERVICE_SERVICE_URL, id);
        logger.info("Sending request to get heating system with id: {}", id);

        HttpHeaders headers = new HttpHeaders();
        headers.setContentType(MediaType.APPLICATION_JSON);
        HttpEntity<Void> requestEntity = new HttpEntity<>(headers);

        try {
            ResponseEntity<HeatingSystemDto> responseEntity = restTemplate.exchange(
                    url,
                    HttpMethod.GET,
                    requestEntity,
                    HeatingSystemDto.class
            );

            if (responseEntity.getStatusCode().is2xxSuccessful()) {
                logger.info("Successfully fetched heating system: {}", responseEntity.getBody());
                return responseEntity.getBody();
            } else {
                logger.error("Failed to fetch heating system: {}", responseEntity.getStatusCode());
                throw new RuntimeException("Failed to fetch heating system. " +
                        "HTTP Status Code: " + responseEntity.getStatusCode());
            }
        } catch (Exception e) {
            logger.error("Error occurred while fetching heating system: {}", e.getMessage(), e);
            throw new RuntimeException("Error occurred while fetching heating system: " + e.getMessage(), e);
        }
    }

    public HeatingSystemDto updateHeatingSystem(Long id, HeatingSystemDto heatingSystemDto) {
        String url = String.format("%s/modules/%d", MICROSERVICE_SERVICE_URL, id);
        logger.info("Sending request to update heating system with id: {}", id);

        HttpHeaders headers = new HttpHeaders();
        headers.setContentType(MediaType.APPLICATION_JSON);
        HttpEntity<HeatingSystemDto> requestEntity = new HttpEntity<>(heatingSystemDto, headers);

        try {
            ResponseEntity<HeatingSystemDto> responseEntity = restTemplate.exchange(
                    url,
                    HttpMethod.PUT,
                    requestEntity,
                    HeatingSystemDto.class
            );

            if (responseEntity.getStatusCode().is2xxSuccessful()) {
                logger.info("Successfully updated heating system: {}", responseEntity.getBody());
                return responseEntity.getBody();
            } else {
                logger.error("Failed to update heating system: {}", responseEntity.getStatusCode());
                throw new RuntimeException("Failed to update heating system. " +
                        "HTTP Status Code: " + responseEntity.getStatusCode());
            }
        } catch (Exception e) {
            logger.error("Error occurred while updating heating system: {}", e.getMessage(), e);
            throw new RuntimeException("Error occurred while updating heating system: " + e.getMessage(), e);
        }
    }

    public HeatingSystemDto sendSetTemperatureRequest(Long moduleId, double temperature) {
        TemperatureRequest temperatureRequest = new TemperatureRequest("set_temperature", temperature);
        String url = String.format("%s/modules/%d/house/%d", MICROSERVICE_SERVICE_URL, moduleId, moduleId);
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

    public HeatingSystemDto getCurrentTemperature(Long moduleId) {
        String url = String.format("%s/modules/%d", MICROSERVICE_SERVICE_URL, moduleId);
        logger.info("Sending request to URL: {}", url);

        HttpHeaders headers = new HttpHeaders();
        // TODO: здесь нужно проставить jwt токен
        headers.setContentType(MediaType.APPLICATION_JSON);
        HttpEntity<Void> requestEntity = new HttpEntity<>(headers);

        try {
            ResponseEntity<HeatingSystemDto> responseEntity = restTemplate.exchange(
                    url,
                    HttpMethod.GET,
                    requestEntity,
                    HeatingSystemDto.class
            );

            if (responseEntity.getStatusCode().is2xxSuccessful()) {
                logger.info("Successfully fetched current temperature: {}", responseEntity.getBody());
                return responseEntity.getBody();
            } else {
                logger.error("Failed to fetch current temperature: {}", responseEntity.getStatusCode());
                throw new RuntimeException("Failed to fetch current temperature. " +
                        "HTTP Status Code: " + responseEntity.getStatusCode());
            }
        } catch (Exception e) {
            logger.error("Error occurred while fetching current temperature: {}", e.getMessage(), e);
            throw new RuntimeException("Error occurred while fetching current temperature: " + e.getMessage(), e);
        }
    }

    public void turnOnHeatingSystem(Long id) {
        String url = String.format("%s/modules/%d/turn-on", MICROSERVICE_SERVICE_URL, id);
        logger.info("Sending request to turn on heating system with id: {}", id);

        HttpHeaders headers = new HttpHeaders();
        headers.setContentType(MediaType.APPLICATION_JSON);
        HttpEntity<Void> requestEntity = new HttpEntity<>(headers);

        try {
            restTemplate.exchange(url, HttpMethod.POST, requestEntity, Void.class);
            logger.info("Successfully turned on heating system with id: {}", id);
        } catch (Exception e) {
            logger.error("Error occurred while turning on heating system: {}", e.getMessage(), e);
            throw new RuntimeException("Error occurred while turning on heating system: " + e.getMessage(), e);
        }
    }

    public void turnOffHeatingSystem(Long id) {
        String url = String.format("%s/modules/%d/turn-off", MICROSERVICE_SERVICE_URL, id);
        logger.info("Sending request to turn off heating system with id: {}", id);

        HttpHeaders headers = new HttpHeaders();
        headers.setContentType(MediaType.APPLICATION_JSON);
        HttpEntity<Void> requestEntity = new HttpEntity<>(headers);

        try {
            restTemplate.exchange(url, HttpMethod.POST, requestEntity, Void.class);
            logger.info("Successfully turned off heating system with id: {}", id);
        } catch (Exception e) {
            logger.error("Error occurred while turning off heating system: {}", e.getMessage(), e);
            throw new RuntimeException("Error occurred while turning off heating system: " + e.getMessage(), e);
        }
    }
}
