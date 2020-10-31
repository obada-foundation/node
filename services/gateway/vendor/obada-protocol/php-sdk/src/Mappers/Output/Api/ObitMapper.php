<?php

declare(strict_types=1);

namespace Obada\Mappers\Output\Api;

use Obada\Mappers\Output\OutputMapper;
use Obada\Obit;
use Obada\Entities\Obit as ApiObit;
use Obada\Properties\Document\Document;
use Obada\Entities\DocumentLink;
use Obada\Properties\Metadata\Record as MetadataRecord;
use Obada\Entities\MetaDataRecord as ApiMetadataRecord;
use Obada\Properties\StructuredData\Record as StructuredDataRecord;
use Obada\Entities\StructureDataRecord as ApiStructureDataRecord;

class ObitMapper implements OutputMapper {
	public function map(Obit $obit) {
		$metadata = array_map(
			fn (MetadataRecord $mdr) =>
				new ApiMetadataRecord(['key' => (string) $mdr->getKey(), 'value' => (string) $mdr->getValue()]),
			$obit->getMetadata()->toArray()
		);

		$structuredData = array_map(
			fn (StructuredDataRecord $sdr) =>
				new ApiStructureDataRecord(['key' => (string) $sdr->getKey(), 'value' => (string) $sdr->getValue()]),
			$obit->getStructuredData()->toArray()
		);

		$documents = array_map(
			fn (Document $d) =>
				new DocumentLink(['name' => (string) $d->getName(), 'hashlink' => (string) $d->getHashLink()]),
			$obit->getDocuments()->toArray()
		);

		return (new ApiObit)
			->setObitDid($obit->getObitId()->toHash())
			->setSerialNumberHash((string) $obit->getSerialNumberHash())
			->setManufacturer((string) $obit->getManufacturer())
			->setPartNumber((string) $obit->getPartNumber())
			->setUsn((string) $obit->getObitId()->toUsn())
			->setObitStatus((string) $obit->getStatus())
			->setOwnerDid((string) $obit->getOwnerDid())
			->setObdDid((string) $obit->getObdDid())
			->setMetadata($metadata)
			->setStructuredData($structuredData)
			->setDocLinks($documents)
			->setModifiedAt((string) $obit->getModifiedAt())// Timestamp is much better here. Discuss with Rohi
			->setRootHash((string) $obit->rootHash());
	}
}