package location_test

import (
	"fmt"
	"testing"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/locationservice"
	"github.com/hashicorp/aws-sdk-go-base/v2/awsv1shim/v2/tfawserr"
	sdkacctest "github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	"github.com/hashicorp/terraform-provider-aws/internal/acctest"
	"github.com/hashicorp/terraform-provider-aws/internal/conns"
)

func TestAccLocationMap_basic(t *testing.T) {
	rName := sdkacctest.RandomWithPrefix(acctest.ResourcePrefix)
	resourceName := "aws_location_map.test"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acctest.PreCheck(t) },
		ErrorCheck:   acctest.ErrorCheck(t, locationservice.EndpointsID),
		Providers:    acctest.Providers,
		CheckDestroy: testAccCheckMapDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccConfigMap_basic(rName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckMapExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "configuration.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "configuration.0.style", "VectorHereBerlin"),
					acctest.CheckResourceAttrRFC3339(resourceName, "create_time"),
					resource.TestCheckResourceAttr(resourceName, "description", ""),
					acctest.CheckResourceAttrRegionalARN(resourceName, "map_arn", "geo", fmt.Sprintf("map/%s", rName)),
					resource.TestCheckResourceAttr(resourceName, "map_name", rName),
					acctest.CheckResourceAttrRFC3339(resourceName, "update_time"),
				),
			},
		},
	})
}

func TestAccLocationMap_description(t *testing.T) {
	rName := sdkacctest.RandomWithPrefix(acctest.ResourcePrefix)
	resourceName := "aws_location_map.test"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acctest.PreCheck(t) },
		ErrorCheck:   acctest.ErrorCheck(t, locationservice.EndpointsID),
		Providers:    acctest.Providers,
		CheckDestroy: testAccCheckMapDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccConfigMap_description(rName, "Test Description"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckMapExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "description", "Test Description"),
				),
			},
		},
	})
}

func testAccCheckMapDestroy(s *terraform.State) error {
	conn := acctest.Provider.Meta().(*conns.AWSClient).LocationConn

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "aws_location_map" {
			continue
		}

		input := &locationservice.DescribeMapInput{
			MapName: aws.String(rs.Primary.ID),
		}

		output, err := conn.DescribeMap(input)

		if tfawserr.ErrCodeEquals(err, locationservice.ErrCodeResourceNotFoundException) {
			continue
		}

		if err != nil {
			return fmt.Errorf("error getting Location Service Map (%s): %w", rs.Primary.ID, err)
		}

		if output != nil {
			return fmt.Errorf("Location Service Map (%s) still exists", rs.Primary.ID)
		}
	}

	return nil
}

func testAccCheckMapExists(resourceName string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[resourceName]

		if !ok {
			return fmt.Errorf("resource not found: %s", resourceName)
		}

		conn := acctest.Provider.Meta().(*conns.AWSClient).LocationConn

		input := &locationservice.DescribeMapInput{
			MapName: aws.String(rs.Primary.ID),
		}

		_, err := conn.DescribeMap(input)

		if err != nil {
			return fmt.Errorf("error getting Location Service Map (%s): %w", rs.Primary.ID, err)
		}

		return nil
	}
}

func testAccConfigMap_basic(rName string) string {
	return fmt.Sprintf(`
resource "aws_location_map" "test" {
  configuration {
    style = "VectorHereBerlin"
  }

  map_name = %[1]q
}
`, rName)
}

func testAccConfigMap_description(rName, description string) string {
	return fmt.Sprintf(`
resource "aws_location_map" "test" {
  configuration {
    style = "VectorHereBerlin"
  }

  map_name    = %[1]q
  description = %[2]q
}
`, rName, description)
}
