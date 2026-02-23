#!/bin/bash
# Interactive script to run scraper job: select container from list, then enter job ID

echo "Running containers:"
echo ""
docker ps --format "table {{.ID}}\t{{.Names}}\t{{.Image}}\t{{.Status}}"
echo ""

containers=()
while IFS= read -r line; do
  containers+=("$line")
done < <(docker ps --format "{{.Names}} ({{.ID}})|{{.ID}}")

if [[ ${#containers[@]} -eq 0 ]]; then
  echo "No running containers found."
  exit 1
fi

echo "Select container (enter number):"
select choice in "${containers[@]}"; do
  if [[ -n "$choice" ]]; then
    container_id="${choice##*|}"
    break
  fi
  echo "Invalid selection. Try again."
done

read -rp "Job ID: " job_id
if [[ -z "$job_id" ]]; then
  echo "Job ID is required."
  exit 1
fi

echo ""
echo "Running: docker exec $container_id ./scraper_app --config scraper_config.yaml --job-id $job_id"
echo ""
docker exec "$container_id" ./scraper_app --config scraper_config.yaml --job-id "$job_id"
